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
	workflowPath := flag.String("f", "", "path to workflow YAML file")
	flag.Parse()
	wf, err := workflow.ParseFile(*workflowPath)
	if err != nil {
		fmt.Print("Parse file error\n")
		os.Exit(1)
	}
	if err := workflow.Validate(wf); err != nil {
		fmt.Print("Validate error\n")
		os.Exit(1)
	}
	ctx := context.Background()
	e := engine.New()

	fmt.Println("Janus Workflow Eninge 0.1.0")
	if err := e.Run(ctx, wf); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
