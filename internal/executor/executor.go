package executor

import (
	"context"

	"github.com/p4sik/janus/internal/workflow"
)

type Executor interface {
	Validate(step workflow.Step) error
	Execute(ctx context.Context, step workflow.Step) error
}
