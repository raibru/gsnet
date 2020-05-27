package pkt

import (
	"bufio"
	"io"
	"os"
	"time"

	"github.com/raibru/gsnet/internal/sys"
)

//
// Logging
//

type pktLogger struct {
	contextName string
}

// LogContext hold logging context
var LogContext = pktLogger{contextName: "pkt"}

// log hold logging context
var logger = sys.ContextLogger{}

func (l pktLogger) ApplyLogger() error {
	err := logger.ApplyLogger(l.contextName)
	if err != nil {
		return err
	}
	logger.Log().Infof("apply package logger behavior: %s", l.contextName)
	logger.Log().Info("::: finish apply package logger")
	return nil
}

func (pktLogger) GetContextName() string {
	return logger.ContextName()
}

//
// Input Packets to send
//

// PacketReader read packet data from input file and distribute it via
// channel which can be used by consumer
type PacketReader struct {
	filename string
	waitMsec time.Duration // wait duration time in milliseconds
	Supply   chan string
}

// NewPacketReader create a new packet reader
func NewPacketReader(name string, wait uint32) *PacketReader {
	m := &PacketReader{
		filename: name,
		waitMsec: time.Duration(wait) * time.Millisecond,
		Supply:   make(chan string),
	}

	return m
}

// Start read packet data
func (pktRead *PacketReader) Start() {
	go func() {
		fn := pktRead.filename

		if pktRead.Supply == nil {
			logger.Log().Errorf("fatal misbehavior data channel is not initialized. Can not provide data from '%s'", fn)
			return
		}

		s, err := os.Stat(fn)
		if err != nil {
			logger.Log().Errorf("failure get os status from. '%s'", fn)
			return
		}

		if s.IsDir() {
			logger.Log().Errorf("failure access input packet file. '%s' is a directory, not a file", fn)
			return
		}

		f, err := os.Open(fn)
		if err != nil {
			logger.Log().Errorf("failure open input packet file '%s'", fn)
			return
		}
		defer f.Close()
		defer close(pktRead.Supply)

		r := bufio.NewReader(f)
		for {
			data, err := r.ReadString('\n')
			if err == io.EOF {
				pktRead.Supply <- "EOF"
				break
			} else if err != nil {
				logger.Log().Errorf("failure read line from input packet file. '%s'", fn)
				pktRead.Supply <- "EOF"
				break
			}
			pktRead.Supply <- data
			time.Sleep(pktRead.waitMsec)
		}
	}()
}

//// Stop read packet data
//func (pktRead *PacketReader) Stop() {
//}
