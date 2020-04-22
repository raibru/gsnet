package sys

import (
	"fmt"
	"os"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

// LoggableContext how to use unified context dependence logger
type LoggableContext interface {
	ApplyLogger() error
	GetContextName() string
}

// LoggingParam hold logging configuration parameter
type LoggingParam struct {
	Service   string `yaml: "service"`
	Version   string `yaml: "version"`
	Filename  string `yaml: "filename"`
	Timestamp string `yaml: "timestamp"`
	Format    string `yaml: "format"`
}

// InitLogging initialize application logging behavior
func InitLogging(lp *LoggingParam) error {
	tsf := lp.Timestamp
	if tsf == "" {
		tsf = "2006-01-02 15:04:05.000"
	}
	// Create the log file if doesn't exist. And append to it if it already exists.
	f, err := os.OpenFile(lp.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	//defer f.Close()

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

	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
		log.SetOutput(os.Stdout)
	} else {
		//log.SetOutput(io.MultiWriter(os.Stdout, f))
		log.SetOutput(f)
	}

	log.Info("==================== Start Logging =====================================")
	log.Infof("run service: %s", lp.Service)
	log.Infof("version    : %s", lp.Version)
	log.Infof("date       : %s", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}

// CreateContextLogging for new Logger with content dependence
func CreateContextLogging(cn string) (*log.Entry, error) {
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
