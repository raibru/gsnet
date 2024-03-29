package cli

import "github.com/raibru/gsnet/cmd/anyserver/etc"

var (
	serverParam etc.ServerServiceParam

	prtVersion     bool
	inputFile      string
	configFile     string
	waitTransfer   uint
	repeatTransfer uint
	teeStdout      bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&prtVersion, "version", "v", false, "Display anyserver version")
	rootCmd.PersistentFlags().StringVarP(&serverParam.Host, "host", "", "", "Listen Tcp/Ip host address")
	rootCmd.PersistentFlags().StringVarP(&serverParam.Port, "port", "", "", "Port of host address to listen")
	rootCmd.PersistentFlags().StringVarP(&serverParam.Name, "name", "", "", "Name of the server service")
	rootCmd.PersistentFlags().StringVarP(&inputFile, "input-file", "i", "", "Use input file send multible data packages")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "f", "", "Use config file for service behavior")
	rootCmd.PersistentFlags().UintVarP(&waitTransfer, "wait-transfer", "", 0, "Wait time to start transfer action in sec")
	rootCmd.PersistentFlags().UintVarP(&repeatTransfer, "repeat-transfer", "", 1, "Repeat transfer action after wait-transfer time")
	rootCmd.PersistentFlags().BoolVarP(&teeStdout, "to-stdout", "", false, "Tee logging output to file and stdout")
}
