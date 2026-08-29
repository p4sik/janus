package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/p4sik/janus/internal/workflow"
)

type CommandExecutor struct{}

type CommandSpec struct {
	Command string `yaml:"command"`
}

func (e *CommandExecutor) Validate(step workflow.Step) error {
	var spec CommandSpec
	if err := step.Spec.Decode(&spec); err != nil {
		return fmt.Errorf("failed to decode command spec: %w", err)
	}

	if spec.Command == "" {
		return fmt.Errorf("command is required")
	}

	return nil
}
func (e *CommandExecutor) Execute(ctx context.Context, step workflow.Step) error {
	var spec CommandSpec
	if err := step.Spec.Decode(&spec); err != nil {
		return fmt.Errorf("failed to decode command spec: %w", err)
	}

	if spec.Command == "" {
		return fmt.Errorf("missing command configuration")
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", spec.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	return nil
}
