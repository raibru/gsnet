package cli

import (
	"fmt"
	"os"

	"github.com/raibru/gsnet/cmd/anyclient/etc"
	"github.com/raibru/gsnet/internal/arch"
	"github.com/raibru/gsnet/internal/pkt"
	"github.com/raibru/gsnet/internal/service"
	"github.com/raibru/gsnet/internal/sys"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "anyclient",
	Short: "AnyClient service",
	Long:  `Use client tcp/ip packet distribution`,
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

	clientService := service.ClientServiceData{}

	if configFile != "" {
		var cf = &etc.AnyClientConfig{}
		err := cf.LoadConfig(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error read config file %s: %s\n", configFile, err.Error())
			os.Exit(2)
		}

		lp := &sys.LoggingParam{
			Service:   cf.Service.Name,
			Version:   Version,
			Filename:  cf.Logging.Filename,
			Timestamp: cf.Logging.Timestamp,
			Format:    cf.Logging.Format,
		}

		if err := sys.InitLogging(lp); err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error initialize logging: %s\n", err.Error())
			os.Exit(2)
		}

		loggables := []sys.LoggableContext{
			sys.LogContext, service.LogContext, pkt.LogContext, arch.LogContext,
		}

		for _, c := range loggables {
			if err := c.ApplyLogger(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: apply logging for: %s -> %s\n", c.GetContextName(), err.Error())
			}
		}

		clientService.Name = cf.Service.Name
		clientService.Addr = cf.Service.Addr
		clientService.Port = cf.Service.Port
	}

	err := clientService.ApplyConnection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Failure. See log. Exit service: %s\n", err.Error())
		os.Exit(2)
	}

	return nil
}
