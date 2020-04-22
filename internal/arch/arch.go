package arch

import (
	"github.com/raibru/gsnet/internal/sys"
	log "github.com/sirupsen/logrus"
)

type archLogger struct {
	doLog       *log.Entry
	contextName string
}

// LogContext hold logging context
var LogContext = archLogger{contextName: "ach"}

func (l archLogger) ApplyLogger() error {
	cl, err := sys.CreateContextLogging(l.contextName)
	if err != nil {
		return err
	}
	l.doLog = cl
	l.doLog.Infof("::: create context logging for: %s", l.contextName)
	return nil
}

func (l archLogger) GetContextName() string {
	return l.contextName
}
