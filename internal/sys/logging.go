package sys

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

//
// Logging configuration
//

var logFile *os.File
var logWriter *bufio.Writer

// LoggingParam hold logging configuration parameter
type LoggingParam struct {
	Service    string
	Version    string
	Filename   string
	TimeFormat string
}

// InitLogging initialize application logging behavior
func InitLogging(lp *LoggingParam) error {
	tsf := lp.TimeFormat
	if tsf == "" {
		tsf = "2006-01-02 15:04:05.000"
	}
	// Create the log file if doesn't exist. And append to it if it already exists.

	log.SetFormatter(&nested.Formatter{

		TimestampFormat: tsf,
		HideKeys:        true,
		NoColors:        true,
		NoFieldsColors:  true,
		ShowFullLevel:   false,
		TrimMessages:    true,
	})
	log.SetLevel(log.TraceLevel)
	log.SetReportCaller(false)

	logFile, err := os.OpenFile(lp.Filename, os.O_WRONLY|os.O_SYNC|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
		log.SetOutput(os.Stdout)
	} else {
		//logWriter := bufio.NewWriter(logFile)
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
		log.RegisterExitHandler(func() {
			fmt.Fprintf(os.Stderr, "... close log file\n")
			log.Info("close log file")
			if logWriter != nil {
				logWriter.Flush()
			}
			if logFile != nil {
				logFile.Close()
			}
		})
		//log.SetOutput(f)
		//defer f.Close()
		//defer w.Flush()
	}

	log.Info("==================== Start Logging =====================================")
	log.Infof("run service: %s", lp.Service)
	log.Infof("version    : %s", lp.Version)
	log.Infof("date       : %s", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}

//
// Logging Context
//

// LoggableContext how to use unified context dependence logger
type LoggableContext interface {
	ApplyLogger() error
	GetContextName() string
	//	Log() *log.Entry
}

// ContextLogger data
type ContextLogger struct {
	contextName string
	logEntry    *log.Entry
}

// ApplyLogger create new named context logger and set ContextLogger data
func (c *ContextLogger) ApplyLogger(cn string) error {
	c.contextName = cn
	e, err := createContextLogging(cn)
	if err != nil {
		return err
	}
	c.logEntry = e
	c.logEntry.Infof("apply context logger for: %s", cn)
	c.logEntry.Tracef("::: create context logging for: %s", cn)
	c.logEntry.Info("::: finish apply context logger")

	return nil
}

// ContextName answer name from ContextLogger data
func (c *ContextLogger) ContextName() string {
	return c.contextName
}

// Log answer log entry object from ContextLogger data
func (c *ContextLogger) Log() *log.Entry {
	return c.logEntry
}

// createContextLogging for new Logger with content dependence
func createContextLogging(cn string) (*log.Entry, error) {
	e := log.WithField("content", "---")
	if cn != "" {
		n := cn
		if len(n) > 3 {
			n = n[:3]
		}
		e = log.WithField("content", n)
	}
	return e, nil
}
