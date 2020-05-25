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
var ctx = sys.ContextLogger{}

func (l archLogger) ApplyLogger() error {
	err := ctx.ApplyLogger(l.contextName)
	if err != nil {
		return err
	}
	ctx.Log().Infof("apply archive logger behavior: %s", l.contextName)
	ctx.Log().Info("::: finish apply archive logger")
	return nil
}

func (archLogger) GetContextName() string {
	return ctx.ContextName()
}

//
// Archive handling
//

var (
	txCount uint32
	rxCount uint32
)

// Record holds send/receive data with meta info per record
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
	Archivate         chan *Record
	filename          string
	archiveType       string
	contextDesciption string
	txCount           uint32
	rxCount           uint32
}

// NewArchive create a new archive object to write archive records
func NewArchive(name string, archType string, ctxDesc string) *Archive {
	a := &Archive{
		Archivate:         make(chan *Record, 10),
		filename:          name,
		archiveType:       archType,
		contextDesciption: ctxDesc,
		txCount:           0,
		rxCount:           0,
	}

	return a
}

// Start starts archiving inside goroutine
func (a *Archive) Start() {
	go func() {
		ctx.Log().Info("start service to write data into archive")
		f, err := os.OpenFile(a.filename, os.O_WRONLY|os.O_SYNC|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			ctx.Log().Errorf("Failure open/create archive file: %s", err.Error())
			return
		}
		defer f.Close()

		w := csv.NewWriter(f)
		defer w.Flush()

		ctx.Log().Info("ready write data into archive")

		for rec := range a.Archivate {

			if rec == nil {
				ctx.Log().Trace("::: receive archive stop event")
				w.Flush()
				return
			}

			ctx.Log().Tracef("::: write data into archive: %s-%d", rec.direction, rec.id)

			data := []string{
				fmt.Sprintf("%s-%d", rec.direction, rec.id),
				rec.time,
				a.contextDesciption,
				rec.direction,
				rec.protocol,
				rec.data}

			if err := w.Write(data); err != nil {
				ctx.Log().Errorf("Failure write data into archive: %s", err.Error())
			}
			w.Flush()
		}
	}()
}

// Stop stops archiving incoming data
func (a *Archive) Stop() {
	a.Archivate <- &Record{}
	ctx.Log().Info("stop service to write data into archive")
}
