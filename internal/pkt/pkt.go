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

// PacketReader read packet data from input file and distribute it via
// channel which can be used by consumer
type PacketReader struct {
	Filename string
	Wait     time.Duration
	Supply   chan string
}

// NewPacketReader create a new packet reader
func NewPacketReader(name string, wait uint32) *PacketReader {
	m := &PacketReader{
		Filename: name,
		Wait:     time.Duration(wait) * time.Millisecond,
		Supply:   make(chan string),
	}

	return m
}

// Start read packet data
func (ctx *PacketReader) Start() {
	go readPacketRawData(ctx)
}

// readPacketRawData read packet data from file which hold raw packet data in each line
func readPacketRawData(pktRead *PacketReader) {
	fn := pktRead.Filename

	if pktRead.Supply == nil {
		ctx.Log().Errorf("fatal misbehavior data channel is not initialized. Can not provide data from '%s'", fn)
		return
	}

	s, err := os.Stat(fn)
	if err != nil {
		ctx.Log().Errorf("failure get os status from. '%s'", fn)
		return
	}

	if s.IsDir() {
		ctx.Log().Errorf("failure access input packet file. '%s' is a directory, not a file", fn)
		return
	}

	f, err := os.Open(fn)
	if err != nil {
		ctx.Log().Errorf("failure open input packet file '%s'", fn)
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		data, err := r.ReadString('\n')
		if err == io.EOF {
			pktRead.Supply <- "EOF"
			break
		} else if err != nil {
			ctx.Log().Errorf("failure read line from input packet file. '%s'", fn)
			pktRead.Supply <- "EOF"
			return //err
		}
		pktRead.Supply <- data
	}

	close(pktRead.Supply)

	return //nil
}
