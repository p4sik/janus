package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/p4sik/janus/internal/engine"
	"github.com/p4sik/janus/internal/executor"
	"github.com/p4sik/janus/internal/validator"
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
		return fmt.Errorf("parse workflow: %w", err)
	}

	executors := map[string]executor.Executor{
		"command":     &executor.CommandExecutor{},
		"http":        &executor.HTTPExecutor{},
		"image_build": &executor.ImageBuildExecutor{},
	}

	if err := validator.Validate(wf, executors); err != nil {
		return fmt.Errorf("validate workflow: %w", err)
	}

	e := engine.New(executors)

	if err := e.Run(context.Background(), wf); err != nil {
		return fmt.Errorf("run workflow: %w", err)
	}

	return nil
}
