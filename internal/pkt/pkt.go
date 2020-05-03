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

//
// Input Packets to send
//

// InputPacketReader hold data about packets read from input file
type InputPacketReader struct {
	Filename string
	WaitSec  time.Duration
	DataChan chan string
}

// NewInputPacketReader create a new input packet reader to read packet data
func NewInputPacketReader(name string, waitSec uint8) *InputPacketReader {
	m := &InputPacketReader{
		Filename: name,
		WaitSec:  5 * time.Second,
		DataChan: make(chan string),
	}

	return m
}

// Start read packet data
func (ctx *InputPacketReader) Start() {
	go readPacketRawData(ctx)
}

// readPacketRawData read packet data from file which hold raw packet data in each line
func readPacketRawData(pktRead *InputPacketReader) {
	fn := pktRead.Filename

	if pktRead.DataChan == nil {
		ctx.Log().Errorf("fatal misbehavior data channel is not initialized. Can not provide data from '%s'", fn)
		return //fmt.Errorf("Packet meta data DataChan shall be initialized")
	}

	s, err := os.Stat(fn)
	if err != nil {
		ctx.Log().Errorf("failure get os status from. '%s'", fn)
		return //err
	}

	if s.IsDir() {
		ctx.Log().Errorf("failure access input packet file. '%s' is a directory, not a file", fn)
		return //fmt.Errorf("'%s' is a directory, not a file", fn)
	}

	f, err := os.Open(fn)
	if err != nil {
		ctx.Log().Errorf("failure open input packet file '%s'", fn)
		return //err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		data, err := r.ReadString('\n')
		if err == io.EOF {
			pktRead.DataChan <- "EOF"
			break
		} else if err != nil {
			ctx.Log().Errorf("failure read line from input packet file. '%s'", fn)
			pktRead.DataChan <- "EOF"
			return //err
		}
		pktRead.DataChan <- data
	}

	close(pktRead.DataChan)

	return //nil
}
