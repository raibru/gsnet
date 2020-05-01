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

// InputPacketContext hold data about packets read from input file
type InputPacketContext struct {
	Filename string
	WaitSec  time.Duration
	DataChan chan string
}

// NewInputPacketContext create a new meta object to read packet data
func NewInputPacketContext(fn string, v uint8) *InputPacketContext {
	m := &InputPacketContext{
		Filename: fn,
		WaitSec:  5 * time.Second,
		DataChan: make(chan string),
	}

	return m
}

// ReadPackets read packet data
func (ctx *InputPacketContext) ReadPackets() {
	go ReadPacketRawData(ctx)
}

// ReadPacketRawData read packet data from file which hold raw packet data in each line
func ReadPacketRawData(pktCtx *InputPacketContext) {
	fn := pktCtx.Filename

	if pktCtx.DataChan == nil {
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
			pktCtx.DataChan <- "EOF"
			break
		} else if err != nil {
			ctx.Log().Errorf("failure read line from input packet file. '%s'", fn)
			pktCtx.DataChan <- "EOF"
			return //err
		}
		pktCtx.DataChan <- data
	}

	close(pktCtx.DataChan)

	return //nil
}
