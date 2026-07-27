package workflow

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Workflow struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

type Step struct {
	ID      string `yaml:"id"`
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Command string `yaml:"command"`
}

func ParseFile(path string) (*Workflow, error) {
	if path == "" {
		return nil, fmt.Errorf("workflow path is required")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow file %q: %w", path, err)
	}

	var wf Workflow

	if err := yaml.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &wf, nil
}

func Validate(wf *Workflow) error {
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

		if step.Type == "command" && step.Command == "" {
			return fmt.Errorf("step %q: command is required", step.ID)
		}
	}
	return nil
}
