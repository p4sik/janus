package executor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/p4sik/janus/internal/workflow"
)

type HTTPExecutor struct{}

type HttpSpec struct {
	URL    string `yaml:"url"`
	Method string `yaml:"method"`
}

func (e *HTTPExecutor) Validate(step workflow.Step) error {
	var spec HttpSpec
	if err := step.Spec.Decode(&spec); err != nil {
		return fmt.Errorf("failed to decode http spec: %w", err)
	}

	if spec.URL == "" {
		return fmt.Errorf("url is required")
	}

	return nil
}

func (e *HTTPExecutor) Execute(ctx context.Context, step workflow.Step) error {
	var spec HttpSpec
	_ = ctx
	if err := step.Spec.Decode(&spec); err != nil {
		return fmt.Errorf("failed to decode http spec: %w", err)
	}

	if spec.URL == "" {
		return fmt.Errorf("missing http configuration")
	}

	req, err := http.NewRequestWithContext(ctx, spec.Method, spec.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute http request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	fmt.Printf(
		"HTTP %s %s -> %s\n", spec.Method, spec.URL, resp.Status)

	return nil
}
