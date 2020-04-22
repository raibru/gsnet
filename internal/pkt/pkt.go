package pkt

import (
	"fmt"
	"os"

	"github.com/raibru/gsnet/internal/sys"
	log "github.com/sirupsen/logrus"
)

type pktLogger struct {
	doLog       *log.Entry
	contextName string
}

var pktLog = pktLogger{contextName: "pkt"}

func (l pktLogger) ApplyLogger() error {
	cl, err := sys.CreateContextLogging(l.contextName)
	if err != nil {
		return err
	}
	l.doLog = cl
	l.doLog.Infof("::: create context logging for: %s", l.contextName)
	return nil
}

// InitPktPackage initialize package behavior
func InitPktPackage() error {
	err := pktLog.ApplyLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error apply logger for content pkt: %s\n", err.Error())
		return err
	}
	return nil
}
