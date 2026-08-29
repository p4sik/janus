package workflow

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

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
