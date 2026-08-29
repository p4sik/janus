package engine

import (
	"context"
	"fmt"

	"github.com/p4sik/janus/internal/executor"
	"github.com/p4sik/janus/internal/workflow"
)

type Engine struct {
	executors map[string]executor.Executor
}

func New(executors map[string]executor.Executor) *Engine {
	return &Engine{executors: executors}
}
func (e *Engine) Run(ctx context.Context, wf *workflow.Workflow) error {
	if wf == nil {
		return fmt.Errorf("workflow is required")
	}

	dependents := make(map[string][]string)
	remainingDependencies := make(map[string]int)

	stepsByID := make(map[string]workflow.Step)

	for _, step := range wf.Steps {
		stepsByID[step.ID] = step
		remainingDependencies[step.ID] = len(step.DependsOn)

		for _, dependencyID := range step.DependsOn {
			dependents[dependencyID] = append(
				dependents[dependencyID],
				step.ID,
			)
		}
	}

	ready := make([]string, 0)

	for stepID, remaining := range remainingDependencies {
		if remaining == 0 {
			ready = append(ready, stepID)
		}
	}

	completed := 0
	for len(ready) > 0 {
		stepID := ready[0]
		ready = ready[1:]

		step := stepsByID[stepID]

		exec, ok := e.executors[step.Type]
		if !ok || exec == nil {
			return fmt.Errorf("Unsuportod executor type: %s", step.Type)
		}

		if err := exec.Execute(ctx, step); err != nil {
			return fmt.Errorf("execute step %s: %w", step.ID, err)
		}

		completed++

		for _, dependentID := range dependents[stepID] {
			remainingDependencies[dependentID]--

			if remainingDependencies[dependentID] == 0 {
				ready = append(ready, dependentID)
			}
		}
	}

	if completed != len(wf.Steps) {
		return fmt.Errorf(
			"workflow cannot be completed: possible dependency cycle",
		)
	}

	_ = ctx

	return nil
}
