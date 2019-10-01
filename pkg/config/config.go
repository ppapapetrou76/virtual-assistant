package config

import (
	"fmt"
	"log"

	"github.com/go-yaml/yaml"
)

// Config is the struct to hold user configuration
type Config struct {
	Labels []string
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
