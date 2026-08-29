package executor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/cli/cli/config"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/p4sik/janus/internal/workflow"
	"github.com/tonistiigi/fsutil"
	"golang.org/x/sync/errgroup"
)

type ImageBuildSpec struct {
	Image      string            `yaml:"image"`
	Context    string            `yaml:"context"`
	Args       map[string]string `yaml:"args"`
	Dockerfile string            `yaml:"dockerfile"`
	Push       bool              `yaml:"push"`
}

type ImageBuildExecutor struct {
}

func (e *ImageBuildExecutor) Validate(step workflow.Step) error {
	spec, err := decodeImageBuildSpec(step)
	if err != nil {
		return err
	}

	if spec.Image == "" {
		return fmt.Errorf("image is required")
	}

	return nil
}

func (e *ImageBuildExecutor) Execute(ctx context.Context, step workflow.Step) error {
	var spec ImageBuildSpec
	if err := step.Spec.Decode(&spec); err != nil {
		return fmt.Errorf("failed to decode build spec: %w", err)
	}
	buildkitAddr := "tcp://127.0.0.1:1234"
	cli, err := client.New(ctx, buildkitAddr)
	if err != nil {
		return fmt.Errorf("failed to create buildkit client: %w", err)
	}
	defer cli.Close()
	contextFS, err := fsutil.NewFS(spec.Context)
	if err != nil {
		return fmt.Errorf("failed to create build context: %w", err)
	}
	dockerfileFS, err := fsutil.NewFS(filepath.Dir(spec.Dockerfile))
	if err != nil {
		return fmt.Errorf("failed to create dockerfile fs: %w", err)
	}
	dockerConfig := config.LoadDefaultConfigFile(os.Stderr)

	authProvider := authprovider.NewDockerAuthProvider(
		authprovider.DockerAuthProviderConfig{
			AuthConfigProvider: authprovider.LoadAuthConfig(dockerConfig),
		},
	)
	solveOpt := client.SolveOpt{
		Frontend: "dockerfile.v0",
		FrontendAttrs: map[string]string{
			"filename": filepath.Base(spec.Dockerfile),
		},
		LocalMounts: map[string]fsutil.FS{
			"context":    contextFS,
			"dockerfile": dockerfileFS,
		},
		Exports: []client.ExportEntry{
			{
				Type: "image",
				Attrs: map[string]string{
					"name": spec.Image,
					"push": fmt.Sprintf("%t", spec.Push),
				},
			},
		},
		Session: []session.Attachable{
			authProvider,
		},
	}
	statusChan := make(chan *client.SolveStatus)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		_, err := cli.Solve(ctx, nil, solveOpt, statusChan)
		if err != nil {
			return fmt.Errorf("failed to solve: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		display, err := progressui.NewDisplay(os.Stdout, progressui.PlainMode)
		if err != nil {
			return fmt.Errorf("failed to create display: %w", err)
		}
		_, err = display.UpdateFrom(ctx, statusChan)
		return err
	})
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}
	return nil
}

func decodeImageBuildSpec(step workflow.Step) (ImageBuildSpec, error) {
	var spec ImageBuildSpec

	if err := step.Spec.Decode(&spec); err != nil {
		return ImageBuildSpec{}, fmt.Errorf("failed to decode build spec: %w", err)
	}

	if spec.Context == "" {
		spec.Context = "."
	}

	if spec.Dockerfile == "" {
		spec.Dockerfile = "Dockerfile"
	}

	return spec, nil
}
