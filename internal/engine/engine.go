package engine

import (
	"context"
	"fmt"

	"github.com/p4sik/janus/internal/workflow"
)

type Engine struct {
}

func New() *Engine {
	return &Engine{}
}
func (e *Engine) Run(ctx context.Context, wf *workflow.Workflow) error {
	if wf == nil {
		return fmt.Errorf("workflow is required")
	}
	_ = ctx
	fmt.Println("Workflow name:", wf.Name)
	names := make([]string, 0, len(wf.Steps))
	for _, step := range wf.Steps {
		names = append(names, step.ID)
	}
	fmt.Println("Steps:", names)
	return nil
}
