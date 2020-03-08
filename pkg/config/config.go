package config

import (
	"fmt"
	"log"

	"github.com/go-yaml/yaml"

	"github.com/ppapapetrou76/virtual-assistant/pkg/util/slices"
)

// Config is the struct to hold user configuration
type Config struct {
	LabelerConfig  `yaml:"labeler"`
	AssignerConfig `yaml:"assigner"`
}

// LabelerConfig is the struct to hold user configuration for the labeler
type LabelerConfig struct {
	IssuesLabelerConfig       `yaml:"issues"`
	PullRequestsLabelerConfig `yaml:"pull-requests"`
}

// IssuesLabelerConfig is the struct to hold user configuration related to issues labeler
type IssuesLabelerConfig struct {
	Labels     slices.StringSlice
	Actions    slices.StringSlice
	OneOfaKind `yaml:"at-least-one"`
}

// OneOfaKind is the struct to hold user configuration related to the feature of checking the existence of at least
// one label of a group and if it doesn't exist then add a default label
type OneOfaKind struct {
	PossibleLabels slices.StringSlice `yaml:"labels"`
	Default        string
}

// PullRequestsLabelerConfig is the struct to hold user configuration related to pull-requests labeler
type PullRequestsLabelerConfig struct {
	Labels  slices.StringSlice
	Actions slices.StringSlice
}

// AssignerConfig is the struct to hold user configuration for the assigner
type AssignerConfig struct {
	IssuesAssignerConfig       `yaml:"issues"`
	PullRequestsAssignerConfig `yaml:"pull-requests"`
}

// PullRequestsAssignerConfig is the struct to hold user configuration related to issues labeler
type PullRequestsAssignerConfig struct {
	Assignee PullRequestsAutoAssigneeConfig `yaml:"assignee"`
	Actions  slices.StringSlice
}

// IssuesAssignerProjectConfig is the struct to hold user configuration related to issues labeler
type PullRequestsAutoAssigneeConfig struct {
	Auto bool `yaml:"auto"`
}

// IssuesAssignerConfig is the struct to hold user configuration related to issues labeler
type IssuesAssignerConfig struct {
	IssuesAssignerProjectConfig `yaml:"project"`
	Actions                     slices.StringSlice
}

// IssuesAssignerProjectConfig is the struct to hold user configuration related to issues labeler
type IssuesAssignerProjectConfig struct {
	ProjectURL string `yaml:"url"`
	Column     string `yaml:"column"`
}

// Load loads config data from raw format to a Config struct
func Load(configRaw *[]byte) (*Config, error) {
	var c = &Config{}

	if configRaw == nil {
		return c, fmt.Errorf("load config : unable to un-marshall empty byte array")
	}

	err := yaml.Unmarshal(*configRaw, c)
	if err != nil {
		return c, fmt.Errorf("load config : unable to un-marshall config [%v], %w", string(*configRaw), err)
	}
	log.Printf("The config: %+v has been successfully unmarshalled", c)

	return c, nil
}
