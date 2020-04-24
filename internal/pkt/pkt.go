package pkt

import "github.com/raibru/gsnet/internal/sys"

type pktLogger struct {
	contextName string
}

// LogContext hold logging context
var LogContext = pktLogger{contextName: "pkt"}

// log hold logging context
var ctx = sys.ContextLogger{}

func (l pktLogger) ApplyLogger() error {
	err := ctx.ApplyLogger(l.contextName)
	if err != nil {
		return err
	}
	ctx.Log().Infof("apply package logger behavior: %s", l.contextName)
	ctx.Log().Info("::: finish apply package logger")
	return nil
}

func (pktLogger) GetContextName() string {
	return ctx.ContextName()
}
