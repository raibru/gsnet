package cli

var prtVersion bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&prtVersion, "version", "v", false, "display gspktservice version")
}
