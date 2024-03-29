package etc

import (
	"fmt"
	"os"

	"github.com/raibru/gsnet/internal/sys"
)

// ClientServiceParam holds data about client services
type ClientServiceParam struct {
	Name       string
	Host       string
	Port       string
	ConfigFile string
}

// AnyClientConfig hold application config environment
type AnyClientConfig struct {
	Service struct {
		Name  string `yaml: "name"`
		Host  string `yaml: "host"`
		Port  string `yaml: "port"`
		Type  string `yaml: "type"`
		Retry uint   `yaml: "retry"`
	} `yaml: "Service"`
	Packet struct {
		Use      bool   `yaml: "use"`
		Filename string `yaml: "filename"`
		Wait     uint   `yaml: "wait"`
	} `yaml: "packet"`
	Archive struct {
		Use      bool   `yaml: "use"`
		Filename string `yaml: "filename"`
		Type     string `yaml: "type"` // yaml, csv, ??? json, xml
	} `yaml: "archive"`
	Logging struct {
		// You can change the Timestamp format. But you have to use the same date and time.
		// "2006-02-02 15:04:06" Works. If you change any digit, it won't work
		Filename   string `yaml: "filename"`
		TimeFormat string `yaml: "timeformat"`
	} `ỳaml: "logging"`
}

// LoadConfig to access given servcie configurations
func (c *AnyClientConfig) LoadConfig(fn string) error {
	if err := sys.LoadConfig(fn, c); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error read config file %s: %s\n", fn, err.Error())
		return err
	}
	return nil
}

// LogFilename answer the filename from config struct where log output will be stored
func LogFilename(c *AnyClientConfig) string {
	return c.Logging.Filename
}
