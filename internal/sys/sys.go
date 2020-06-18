package sys

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// Applier apply system behavior
type Applier interface {
	Apply() error
}

// Identifier identify used logger
type Identifier interface {
	Identify() string
}

type sysLogger struct {
	contextName string
}

// LogContext hold logging context
var LogContext = sysLogger{contextName: "sys"}

// log hold logging context
var logger = LoggerEntity{}

func (l sysLogger) Apply() error {
	err := logger.ApplyLogger(l.contextName)
	if err != nil {
		return err
	}
	logger.Log().Infof("apply system logger behavior: %s", l.contextName)
	logger.Log().Info("finish apply system logger")
	return nil
}

func (sysLogger) Identify() string {
	return logger.ContextName()
}

// Exit stop the logging and exit
func Exit(l int) {
	// execute registerd logging exit handler
	log.Exit(l)
}

// StartSignalHandler run a system signal 'listener' on a new goroutine which will notify by os interrupts
func StartSignalHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Log().Info("handle system signal SIGTERM")
		Exit(0) // use own exit func
	}()
}

// validateFileExists just makes sure, that the path provided is a file,
func validateFileExists(fn string) error {
	s, err := os.Stat(fn)
	if err != nil {
		return err
	}
	if s.IsDir() {
		logger.Log().Errorf("failure access logging file. '%s' is a directory, not a file", fn)
		return fmt.Errorf("'%s' is a directory, not a file", fn)
	}
	return nil
}
