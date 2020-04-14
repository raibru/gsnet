package cli

import "github.com/raibru/gsnet/internal/service"

var (
	serverService service.ServerServiceData

	prtVersion bool
	inputfile  string
	configfile string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&prtVersion, "version", "v", false, "Display anyserver version")
	rootCmd.PersistentFlags().StringVarP(&serverService.Addr, "address", "", "", "Listen Tcp/Ip address")
	rootCmd.PersistentFlags().StringVarP(&serverService.Port, "port", "", "", "Port of address")
	rootCmd.PersistentFlags().StringVarP(&serverService.Name, "name", "", "", "Name of the server service")
	rootCmd.PersistentFlags().StringVarP(&inputfile, "input-file", "i", "", "Use input file send multible data packages")
	rootCmd.PersistentFlags().StringVarP(&configfile, "config-file", "f", "", "Use config file for service behavior")
}
