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
	ctx.Log().Infof("apply system logger behavior: %s", l.contextName)
	ctx.Log().Info("::: finish apply system logger")
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
		ctx.Log().Errorf("failure access logging file. '%s' is a directory, not a file", fn)
		return fmt.Errorf("'%s' is a directory, not a file", fn)
	}
	return nil
}
