package arch

import "github.com/raibru/gsnet/internal/sys"

type archLogger struct {
	contextName string
}

// LogContext hold logging context
var LogContext = archLogger{contextName: "ach"}

// log hold logging context
var ctx = sys.ContextLogger{}

func (l archLogger) ApplyLogger() error {
	err := ctx.ApplyLogger(l.contextName)
	if err != nil {
		return err
	}
	ctx.Log().Infof("::: use package 'arch' wide logging with context: %s", l.contextName)
	return nil
}

func (archLogger) GetContextName() string {
	return ctx.ContextName()
}
