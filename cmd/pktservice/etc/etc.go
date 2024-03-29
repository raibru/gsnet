package etc

import (
	"fmt"
	"os"

	"github.com/raibru/gsnet/internal/sys"
)

// PktServiceParam holds data about server services
type PktServiceParam struct {
	Name       string
	Host       string
	Port       string
	ConfigFile string
}

// PktServiceConfig hold application config environment
type PktServiceConfig struct {
	Service struct {
		Name    string `yaml: "name"`
		Network []struct {
			Channel struct {
				Name          string `yaml: "name"`
				Type          string `yaml: "type"`
				ReconInterval uint32 `yaml: "recon_interval"`
				Listener      struct {
					Name string `yaml: "name"`
					Host string `yaml: "host"`
					Port string `yaml: "port"`
				} `yaml: "listener"`
				Dialer struct {
					Name  string `yaml: "name"`
					Host  string `yaml: "host"`
					Port  string `yaml: "port"`
					Retry uint   `yaml: "retry"`
				} `yaml: "dialer"`
			} `yaml: "channel,flow"`
		} `yaml: "network"`
	} `yaml: "service"`
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
func (c *PktServiceConfig) LoadConfig(fn string) error {
	if err := sys.LoadConfig(fn, c); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error read config file %s: %s\n", fn, err.Error())
		return err
	}
	return nil
}

// LogFilename answer the filename from config struct where log output will be stored
func LogFilename(c *PktServiceConfig) string {
	return c.Logging.Filename
}
