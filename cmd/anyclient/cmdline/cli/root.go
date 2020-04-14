package cli

import "github.com/raibru/gsnet/internal/service"

var (
	clientService service.ClientServiceData

	prtVersion bool
	inputfile  string
	configfile string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&prtVersion, "version", "v", false, "display anxclient version")
	rootCmd.PersistentFlags().StringVarP(&clientService.Addr, "address", "", "", "connect to Tcp/Ip address")
	rootCmd.PersistentFlags().StringVarP(&clientService.Port, "port", "", "", "Port of address")
	rootCmd.PersistentFlags().StringVarP(&clientService.Name, "name", "", "", "Name of the client service")
	rootCmd.PersistentFlags().StringVarP(&inputfile, "input-file", "i", "", "Use input file send multible data packages")
	rootCmd.PersistentFlags().StringVarP(&configfile, "config-file", "f", "", "Use config file for service behavior")
}
