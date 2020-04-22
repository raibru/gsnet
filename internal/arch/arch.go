package arch

import (
	"fmt"
	"os"

	"github.com/raibru/gsnet/internal/sys"
	log "github.com/sirupsen/logrus"
)

type archLogger struct {
	doLog       *log.Entry
	contextName string
}

var archLog = archLogger{contextName: "ach"}

func (l archLogger) ApplyLogger() error {
	cl, err := sys.CreateContextLogging(l.contextName)
	if err != nil {
		return err
	}
	l.doLog = cl
	l.doLog.Infof("::: create context logging for: %s", l.contextName)
	return nil
}

// InitArchPackage initialize package behavior
func InitArchPackage() error {
	err := archLog.ApplyLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error apply logger for content arch: %s\n", err.Error())
		return err
	}
	return nil
}
