package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/p4sik/janus/internal/engine"
	"github.com/p4sik/janus/internal/workflow"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	workflowPath := flag.String("f", "", "path to workflow YAML file")
	flag.Parse()

	wf, err := workflow.ParseFile(*workflowPath)
	if err != nil {
		return err
	}

	if err := workflow.Validate(wf); err != nil {
		return fmt.Errorf("validate workflow: %w", err)
	}

	e := engine.New()

	if err := e.Run(context.Background(), wf); err != nil {
		return fmt.Errorf("run workflow: %w", err)
	}

	return nil
}
