package pkt

import (
	"github.com/raibru/gsnet/internal/sys"
	log "github.com/sirupsen/logrus"
)

type pktLogger struct {
	doLog       *log.Entry
	contextName string
}

// LogContext hold logging context
var LogContext = pktLogger{contextName: "pkt"}

func (l pktLogger) ApplyLogger() error {
	cl, err := sys.CreateContextLogging(l.contextName)
	if err != nil {
		return err
	}
	l.doLog = cl
	l.doLog.Infof("::: create context logging for: %s", l.contextName)
	return nil
}

func (l pktLogger) GetContextName() string {
	return l.contextName
}
