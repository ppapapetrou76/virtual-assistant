package config

import (
	"fmt"
	"log"

	"github.com/go-yaml/yaml"

	"github.com/ppapapetrou76/virtual-assistant/pkg/util/slices"
)

// Config is the struct to hold user configuration
type Config struct {
	IssuesConfig       `yaml:"issues"`
	PullRequestsConfig `yaml:"pull-requests"`
}

// IssuesConfig is the struct to hold user configuration related to issues
type IssuesConfig struct {
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

// PullRequestsConfig is the struct to hold user configuration related to pull-requests
type PullRequestsConfig struct {
	Labels  slices.StringSlice
	Actions slices.StringSlice
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
