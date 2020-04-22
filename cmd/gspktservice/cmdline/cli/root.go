package cli

import "github.com/raibru/gsnet/cmd/gspktservice/etc"

var (
	gspktParam etc.GsPktServiceParam

	prtVersion bool
	inputFile  string
	configFile string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&prtVersion, "version", "v", false, "Display anyserver version")
	rootCmd.PersistentFlags().StringVarP(&gspktParam.Name, "name", "", "", "Name of the server service")
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input-file", "i", "", "Use input file send multible data packages")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "f", "", "Use config file for service behavior")
}
