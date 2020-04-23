package sys

import (
	"fmt"
	"os"
)

type sysLogger struct {
	contextName string
}

// LogContext hold logging context
var LogContext = sysLogger{contextName: "sys"}

// log hold logging context
var ctx = ContextLogger{}

func (l sysLogger) ApplyLogger() error {
	err := ctx.ApplyLogger(l.contextName)
	if err != nil {
		return err
	}
	ctx.Log().Infof("::: use package 'sys' wide logging with context: %s", l.contextName)
	return nil
}

func (sysLogger) GetContextName() string {
	return ctx.ContextName()
}

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
