package cli

import "github.com/raibru/gsnet/cmd/anyclient/etc"

var (
	clientParam etc.ClientServiceParam

	prtVersion bool
	inputFile  string
	configFile string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&prtVersion, "version", "v", false, "display anxclient version")
	rootCmd.PersistentFlags().StringVarP(&clientParam.Addr, "address", "", "", "connect to Tcp/Ip address")
	rootCmd.PersistentFlags().StringVarP(&clientParam.Port, "port", "", "", "Port of address")
	rootCmd.PersistentFlags().StringVarP(&clientParam.Name, "name", "", "", "Name of the client service")
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input-file", "i", "", "Use input file send multible data packages")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "f", "", "Use config file for service behavior")
}
