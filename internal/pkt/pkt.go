package pkt

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
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
var logger = sys.LoggerEntity{}

func (l pktLogger) Apply() error {
	err := logger.ApplyLogger(l.contextName)
	if err != nil {
		return err
	}
	logger.Log().Infof("apply package logger behavior: %s", l.contextName)
	logger.Log().Info("finish apply package logger")
	return nil
}

func (pktLogger) Identify() string {
	return logger.ContextName()
}

//
// Input Packets to supply
//

// PacketReader read packet data from input file and distribute it via
// channel which can be used by consumer
type PacketReader struct {
	filename     string
	waitMsec     time.Duration // wait duration time in milliseconds
	waitStartSec time.Duration // wait duration until supply data is started
	Use          bool
	Supply       chan []byte
}

// SetSupply set supply channel for data which have to process
func (s *PacketReader) SetSupply(c chan []byte) {
	s.Supply = c
}

// NewPacketReader create a new packet reader
func NewPacketReader(name string, wait uint, waitStart uint) *PacketReader {
	return &PacketReader{
		filename:     name,
		waitMsec:     time.Duration(wait) * time.Millisecond,
		waitStartSec: time.Duration(waitStart) * time.Second,
		Use:          true,
		Supply:       nil,
	}
}

// NonPacketReader create a new packet reader
func NonPacketReader() *PacketReader {
	return &PacketReader{
		Use: false,
	}
}

// Start read packet data
func (pktRead *PacketReader) Start(done chan bool) {
	time.Sleep(pktRead.waitStartSec)
	go handle(pktRead, done)
}

func handle(pktRead *PacketReader, done chan bool) {
	fn := pktRead.filename

	if pktRead.Supply == nil {
		logger.Log().Errorf("fatal misbehavior data channel is not initialized. Can not provide data from '%s'", fn)
		done <- true
		return
	}

	s, err := os.Stat(fn)
	if err != nil {
		logger.Log().Errorf("failure get os status from. '%s'", fn)
		done <- true
		return
	}

	if s.IsDir() {
		logger.Log().Errorf("failure access input packet file. '%s' is a directory, not a file", fn)
		done <- true
		return
	}

	f, err := os.Open(fn)
	if err != nil {
		logger.Log().Errorf("failure open input packet file '%s'", fn)
		done <- true
		return
	}
	defer f.Close()
	defer close(pktRead.Supply)

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			logger.Log().Trace("read EOF from data file")
			break
		} else if err != nil {
			logger.Log().Errorf("failure read line from input packet file. '%s'", fn)
			break
		}

		if match, _ := regexp.MatchString(`^#`, line); match {
			continue
		} else if match, _ := regexp.MatchString(`^\w*$`, line); match {
			continue
		} else if match, _ := regexp.MatchString(`^EOF`, line); match {
			break
		}

		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, "\r", "", -1)
		if len(line) > 0 {
			pktRead.Supply <- []byte(line)
			time.Sleep(pktRead.waitMsec)
		}
	}
	//pktRead.Supply <- "EOF"
	done <- true
}

//// Stop read packet data
//func (pktRead *PacketReader) Stop() {
//}
