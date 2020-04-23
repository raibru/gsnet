package sys

import (
	"os"

	"gopkg.in/yaml.v2"
)

// ServiceConfig defines the behavior how to handle
// configuration for all services
type ServiceConfig interface {
	LoadConfig(fn string) error
}

// LoadConfig set configuration parameter from given filename fn
func LoadConfig(fn string, conf ServiceConfig) error {
	if verr := validateFileExists(fn); verr != nil {
		return verr
	}

	// Open config file
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(f)

	// Start YAML decoding from file
	if err := d.Decode(conf); err != nil {
		return err
	}

	return nil
}