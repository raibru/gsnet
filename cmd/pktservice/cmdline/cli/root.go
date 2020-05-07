package cli

import "github.com/raibru/gsnet/cmd/pktservice/etc"

var (
	pktParam etc.PktServiceParam

	prtVersion bool
	inputFile  string
	configFile string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&prtVersion, "version", "v", false, "Display anyserver version")
	rootCmd.PersistentFlags().StringVarP(&pktParam.Name, "name", "", "", "Name of the server service")
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input-file", "i", "", "Use input file send multible data packages")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "f", "", "Use config file for service behavior")
}
