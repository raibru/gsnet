package sys

import (
	"fmt"
	"os"
)

// validateFileExists just makes sure, that the path provided is a file,
func validateFileExists(fn string) error {
	s, err := os.Stat(fn)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", fn)
	}
	return nil
}
