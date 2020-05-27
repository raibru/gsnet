package sys

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ServiceConfig defines the behavior how to handle
// configuration for all services
type ServiceConfig interface {
	LoadConfig(fn string) error
}

// LoadConfig set configuration parameter from given filename fn
func LoadConfig(fn string, conf ServiceConfig) error {
	if err := validateFileExists(fn); err != nil {
		return err
	}

	bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		return err
	}

	return nil
}
