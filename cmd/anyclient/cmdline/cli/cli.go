package cli

import (
	"fmt"
	"os"

	"github.com/raibru/gsnet/cmd/anyclient/etc"
	"github.com/raibru/gsnet/internal/archive"
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

	sys.StartSignalHandler()
	var clientService *service.ClientServiceData
	var archiveService *archive.Archive
	var readerService *pkt.PacketReader

	if configFile != "" {
		var cf = &etc.AnyClientConfig{}
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
		readerService = pkt.NewPacketReader(cf.Packet.Filename, cf.Packet.Wait)
		clientService = service.NewClientService(
			cf.Service.Name,
			cf.Service.Host,
			cf.Service.Port,
			readerService.Supply,
			archiveService.Archivate)
	} else {
		clientService = service.NewClientService(
			"anyclient",
			"127.0.0.1",
			"30100",
			nil,
			nil)
	}

	err := clientService.ApplyConnection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Failure. See log. Exit service: %s\n", err.Error())
		sys.Exit(2)
	}

	defer clientService.Finalize()

	wait := make(chan bool, 1)
	readed := make(chan bool, 1)
	sent := make(chan bool, 1)

	if archiveService != nil {
		archiveService.Start(wait)
	}
	if readerService != nil {
		readerService.Start(readed)
	}
	clientService.SendPackets(sent)

	<-sent

	return nil
}
