package cli

import (
	"fmt"
	"os"

	"github.com/raibru/gsnet/cmd/gspktservice/etc"
	"github.com/raibru/gsnet/internal/sys"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gspktservice",
	Short: "Ground Station Packet Service",
	Long:  `Provide ground station packet service behavior via tcp/ip comminication`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleParam(cmd, args); err != nil {
			cmd.Help()
			fmt.Println("\nRoot command has parsing error: ", err)
			os.Exit(1)
		}
	},
}

// Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Execute root cmd has error", err)
		os.Exit(1)
	}
}

// handleParam parameter evaluation
func handleParam(cmd *cobra.Command, args []string) error {
	if prtVersion {
		PrintVersion(os.Stdout)
		return nil
	}

	if configFile != "" {
		var cf = &etc.GsPktServiceConfig{}
		err := cf.LoadConfig(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error read config file %s: %s\n", configFile, err.Error())
			os.Exit(2)
		}

		lp := &sys.LoggingParam{
			Filename:  cf.Logging.Filename,
			Timestamp: cf.Logging.Timestamp,
			Format:    cf.Logging.Format,
		}

		if err := sys.InitLogging(lp); err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error initialize logging: %s\n", err.Error())
			os.Exit(2)
		}
	}

	return nil
}
