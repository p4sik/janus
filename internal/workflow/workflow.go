package workflow

import (
	"gopkg.in/yaml.v3"
)

type Workflow struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

type Step struct {
	ID        string    `yaml:"id"`
	Name      string    `yaml:"name"`
	Type      string    `yaml:"type"`
	Spec      yaml.Node `yaml:"spec,omitempty"`
	DependsOn []string  `yaml:"dependsOn"`
}
