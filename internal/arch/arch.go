package arch

import (
	"encoding/csv"
	"fmt"
	"os"

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

// ArchiveRecord holds send/receive data with meta info per record
type ArchiveRecord struct {
	MsgID        uint32
	MsgTime      string
	MsgDirection string // RX, TX
	Protocol     string
	Data         string
}

// Archive hold archive runable parameter
type Archive struct {
	Filename    string
	ArchiveType string
	DataChan    chan ArchiveRecord
	ServName    string
	TxCount     uint32
	RxCount     uint32
}

// NewArchive create a new archive object to write archive records
func NewArchive(name string, archType string, servName string) *Archive {
	a := &Archive{
		Filename:    name,
		ArchiveType: archType,
		ServName:    servName,
		TxCount:     0,
		RxCount:     0,
		DataChan:    make(chan ArchiveRecord, 10),
	}

	return a
}

// Start run archiving in goroutine
func (a *Archive) Start() {
	go handleArchive(a)
}

// handleArchive archive data into configured archive destination
func handleArchive(a *Archive) {
	ctx.Log().Info("handle archive data")
	f, err := os.OpenFile(a.Filename, os.O_WRONLY|os.O_SYNC|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		ctx.Log().Errorf("Failure open/create archive file: %s", err.Error())
		return
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	ctx.Log().Info("ready write data into archive")

	for rec := range a.DataChan {

		ctx.Log().Tracef("::: write data into archive: %d", rec.MsgID)

		data := []string{
			fmt.Sprint(rec.MsgID),
			rec.MsgTime,
			a.ServName,
			rec.MsgDirection,
			rec.Protocol,
			rec.Data}

		if err := w.Write(data); err != nil {
			ctx.Log().Errorf("Failure write data into archive: %s", err.Error())
		}
		w.Flush()
	}
	return
}
