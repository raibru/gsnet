package etc

import (
	"fmt"
	"os"

	"github.com/raibru/gsnet/internal/sys"
)

// ClientServiceParam holds data about client services
type ClientServiceParam struct {
	Name       string
	Addr       string
	Port       string
	ConfigFile string
}

// AnyClientConfig hold application config environment
type AnyClientConfig struct {
	Service struct {
		ServiceName string `yaml: "serviceName"`
		IPAddr      string `yaml: "ipAddr"`
		Port        string `yaml: "port"`
	} `yaml: "Service"`
	Logging struct {
		// You can change the Timestamp format. But you have to use the same date and time.
		// "2006-02-02 15:04:06" Works. If you change any digit, it won't work
		Filename  string `yaml: "filename"`
		Timestamp string `yaml: "timestamp"`
		Format    string `yaml: "format"`
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
