package etc

import (
	"fmt"
	"os"

	"github.com/raibru/gsnet/internal/sys"
)

// ServerServiceParam holds data about server services
type ServerServiceParam struct {
	Name       string
	Addr       string
	Port       string
	ConfigFile string
}

// AnyServerConfig hold application config environment
type AnyServerConfig struct {
	Service struct {
		Name string `yaml: "name"`
		Addr string `yaml: "addr"`
		Port string `yaml: "port"`
	} `yaml: "Service"`
	Packet struct {
		Filename string `yaml: "filename"`
		Wait     uint   `yaml:"wait"`
	} `yaml: "Packet"`
	Archive struct {
		Filename string `yaml: "filename"`
		Type     string `yaml: "type"` // yaml, csv, ??? json, xml
	} `yaml: "archive"`
	Logging struct {
		// You can change the Timestamp format. But you have to use the same date and time.
		// "2006-02-02 15:04:06" Works. If you change any digit, it won't work
		Filename  string `yaml: "filename"`
		Timestamp string `yaml: "timestamp"`
		Format    string `yaml: "format"`
	} `ỳaml: "logging"`
}

// LoadConfig to access given servcie configurations
func (c *AnyServerConfig) LoadConfig(fn string) error {
	if err := sys.LoadConfig(fn, c); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error read config file %s: %s\n", fn, err.Error())
		return err
	}
	return nil
}

// LogFilename answer the filename from config struct where log output will be stored
func LogFilename(c *AnyServerConfig) string {
	return c.Logging.Filename
}
