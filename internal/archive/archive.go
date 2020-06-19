package archive

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/raibru/gsnet/internal/sys"
)

//
// Context logging
//

type archLogger struct {
	contextName string
}

// LogContext hold logging context
var LogContext = archLogger{contextName: "ach"}

// log hold logging context
var logger = sys.LoggerEntity{}

func (l archLogger) Apply() error {
	err := logger.ApplyLogger(l.contextName)
	if err != nil {
		return err
	}
	logger.Log().Infof("apply archive logger behavior: %s", l.contextName)
	logger.Log().Info("finish apply archive logger")
	return nil
}

func (archLogger) Identify() string {
	return logger.ContextName()
}

//
// Archive handling
//

var (
	txCount uint32
	rxCount uint32
)

// Record holds transfered/received data with meta info per record
type Record struct {
	id        uint32
	time      string
	direction string // RX, TX
	protocol  string
	data      string
}

// NewRecord create an archive record with time and message id
func NewRecord(hexData string, rxtx string, protocol string) *Record {
	var count uint32
	switch rxtx {
	case "TX":
		txCount++
		count = txCount
	case "RX":
		rxCount++
		count = rxCount
	default:
		count = 0
	}
	t := time.Now().Format("2006-01-02 15:04:05.000")
	r := &Record{
		id:        count,
		time:      t,
		direction: rxtx,
		protocol:  protocol,
		data:      hexData}
	return r
}

// Archive hold archive runable parameter
type Archive struct {
	Use               bool
	filename          string
	archiveType       string
	contextDesciption string
	txCount           uint32
	rxCount           uint32
	Archivate         chan *Record
}

// NewArchive create a new archive object to write archive records
func NewArchive(name string, archType string, ctxDesc string) *Archive {
	return &Archive{
		Use:               true,
		filename:          name,
		archiveType:       archType,
		contextDesciption: ctxDesc,
		txCount:           0,
		rxCount:           0,
		Archivate:         nil,
		//Archivate:         make(chan *Record, 10),
	}
}

// NonArchive create a new archive object to write archive records
func NonArchive() *Archive {
	return &Archive{
		Use: false,
	}
}

// SetArchivate set process data channel
func (a *Archive) SetArchivate(c chan *Record) {
	a.Archivate = c
}

// Start starts archiving inside goroutine
func (a *Archive) Start(done chan bool) {
	go handle(a, done)
}

func handle(a *Archive, done chan bool) {
	logger.Log().Info("start service to write data into archive")
	f, err := os.OpenFile(a.filename, os.O_WRONLY|os.O_SYNC|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logger.Log().Errorf("Failure open/create archive file: %s", err.Error())
		return
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	done = make(chan bool)
	logger.Log().Info("ready write data into archive")

	for {
		rec, more := <-a.Archivate

		if !more {
			logger.Log().Trace("receive archive stop event")
			w.Flush()
			done <- true
			return
		}
		//		if rec == nil {
		//			logger.Log().Trace("receive archive stop event")
		//			w.Flush()
		//			return
		//		}

		logger.Log().Tracef("write data into archive: %s-%d", rec.direction, rec.id)

		data := []string{
			fmt.Sprintf("%s-%d", rec.direction, rec.id),
			rec.time,
			a.contextDesciption,
			rec.direction,
			rec.protocol,
			rec.data}

		if err := w.Write(data); err != nil {
			logger.Log().Errorf("Failure write data into archive: %s", err.Error())
		}
		w.Flush()
	}
}

// Stop stops archiving incoming data
func (a *Archive) Stop() {
	//a.Archivate <- &Record{}
	close(a.Archivate)
	logger.Log().Info("stop service write data into archive")
}
