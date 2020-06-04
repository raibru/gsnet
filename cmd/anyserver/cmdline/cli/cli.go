package cli

import (
	"fmt"
	"os"

	"github.com/raibru/gsnet/cmd/anyserver/etc"
	"github.com/raibru/gsnet/internal/archive"
	"github.com/raibru/gsnet/internal/pkt"
	"github.com/raibru/gsnet/internal/service"
	"github.com/raibru/gsnet/internal/sys"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "anyserver",
	Short: "AnyServer service",
	Long:  `Provide server tcp/ip packet exchange`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleParam(cmd, args); err != nil {
			cmd.Help()
			fmt.Println("\nRoot command has parsing error: ", err)
			sys.Exit(1)
		}
	},
}

// Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Execute root cmd has error", err)
		sys.Exit(1)
	}
}

// handleParam parameter evaluation
func handleParam(cmd *cobra.Command, args []string) error {
	if prtVersion {
		PrintVersion(os.Stdout)
		return nil
	}

	var srvService *service.ServerServiceData
	var archiveService *archive.Archive

	sys.StartSignalHandler()

	if configFile != "" {
		var cf = etc.AnyServerConfig{}
		err := cf.LoadConfig(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error read config file %s: %s\n", configFile, err.Error())
			sys.Exit(2)
		}

		lp := &sys.LoggingParam{
			Service:    cf.Service.Name,
			Version:    Version,
			Filename:   cf.Logging.Filename,
			TimeFormat: cf.Logging.TimeFormat,
		}

		if err := sys.InitLogging(lp); err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error initialize logging: %s\n", err.Error())
			sys.Exit(2)
		}

		loggables := []sys.ContextLogger{
			sys.LogContext, service.LogContext, pkt.LogContext, archive.LogContext,
		}

		for _, l := range loggables {
			if err := l.Apply(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: apply logging for: %s -> %s\n", l.Identify(), err.Error())
			}
		}

		archiveService = archive.NewArchive(cf.Archive.Filename, cf.Archive.Type, cf.Service.Name)
		srvService = service.NewServerService(
			cf.Service.Name,
			cf.Service.Host,
			cf.Service.Port,
			nil,
			archiveService.Archivate)
	} else {
		srvService = service.NewServerService(
			"anyserver",
			"127.0.0.1",
			"30100",
			nil,
			nil)
	}

	wait := make(chan bool)

	if archiveService != nil {
		archiveService.Start(wait)
	}

	err := srvService.ApplyConnection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Failure. See log. Exit service: %s\n", err.Error())
		sys.Exit(2)
	}

	return nil
}
