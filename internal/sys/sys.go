package sys

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type sysLogger struct {
	doLog       *log.Entry
	contextName string
}

var sysLog = sysLogger{contextName: "sys"}

func (l sysLogger) ApplyLogger() error {
	cl, err := CreateContextLogging(l.contextName)
	if err != nil {
		return err
	}
	l.doLog = cl
	l.doLog.Infof("::: create context logging for: %s", l.contextName)
	return nil
}

// InitSysPackage initialize package behavior
func InitSysPackage() error {
	err := sysLog.ApplyLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error apply logger for content sys: %s\n", err.Error())
		return err
	}
	return nil
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
