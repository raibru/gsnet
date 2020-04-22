package sys

import (
	"fmt"
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

// LoggableCategory how to use logger
type LoggableCategory interface {
	ApplyLogger(cat string) error
}

type sysLogger struct{ doLog *log.Entry }

var sysLog = sysLogger{}

func (s sysLogger) ApplyLogger(cn string) error {
	cat, err := CreateCategoryLogging(cn)
	if err != nil {
		return err
	}
	s.doLog = cat
	s.doLog.Infof("::: create and apply logging for category %s", cn)
	return nil
}

// InitSysPackage initialize package behavior
func InitSysPackage() error {
	err := sysLog.ApplyLogger("sys")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error apply logger for category sys: %s\n", err.Error())
		return err
	}
	return nil
}

// LoggingParam hold loggin configuration parameter
type LoggingParam struct {
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

	return nil
}

// CreateCategoryLogging for new Logger with category
func CreateCategoryLogging(cn string) (*log.Entry, error) {
	e := log.WithField("category", cn)
	return e, nil
}
