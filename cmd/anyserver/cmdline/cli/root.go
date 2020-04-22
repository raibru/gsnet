package cli

import "github.com/raibru/gsnet/cmd/anyserver/etc"

var (
	serverParam etc.ServerServiceParam

	prtVersion bool
	inputFile  string
	configFile string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&prtVersion, "version", "v", false, "Display anyserver version")
	rootCmd.PersistentFlags().StringVarP(&serverParam.Addr, "address", "", "", "Listen Tcp/Ip address")
	rootCmd.PersistentFlags().StringVarP(&serverParam.Port, "port", "", "", "Port of address")
	rootCmd.PersistentFlags().StringVarP(&serverParam.Name, "name", "", "", "Name of the server service")
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input-file", "i", "", "Use input file send multible data packages")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "f", "", "Use config file for service behavior")
}
