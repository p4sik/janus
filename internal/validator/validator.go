package validator

import (
	"fmt"

	"github.com/p4sik/janus/internal/executor"
	"github.com/p4sik/janus/internal/workflow"
)

func Validate(wf *workflow.Workflow, executors map[string]executor.Executor) error {
	if wf == nil {
		return fmt.Errorf("workflow is required")
	}

	if wf.Name == "" {
		return fmt.Errorf("workflow name is required")
	}

	if len(wf.Steps) == 0 {
		return fmt.Errorf("workflow must contain at least one step")
	}

	seenIDs := make(map[string]struct{}, len(wf.Steps))

	for i, step := range wf.Steps {
		if step.ID == "" {
			return fmt.Errorf("step at index %d: id is required", i)
		}

		if _, exists := seenIDs[step.ID]; exists {
			return fmt.Errorf("duplicate step id %q", step.ID)
		}

		seenIDs[step.ID] = struct{}{}

		if step.Type == "" {
			return fmt.Errorf("step %q: type is required", step.ID)
		}

		stepExecutor, ok := executors[step.Type]
		if !ok || stepExecutor == nil {
			return fmt.Errorf("step %q: unsupported type %q", step.ID, step.Type)
		}

		if err := stepExecutor.Validate(step); err != nil {
			return fmt.Errorf("step %q: invalid %s spec: %w", step.ID, step.Type, err)
		}

		for _, dependencyID := range step.DependsOn {
			if dependencyID == "" {
				return fmt.Errorf("step %q: dependency id is required", step.ID)
			}

			if dependencyID == step.ID {
				return fmt.Errorf("step %q: cannot depend on itself", step.ID)
			}
		}
	}

	for _, step := range wf.Steps {
		for _, dependencyID := range step.DependsOn {
			if _, exists := seenIDs[dependencyID]; !exists {
				return fmt.Errorf("step %q: unknown dependency %q", step.ID, dependencyID)
			}
		}
	}

	return nil
}
