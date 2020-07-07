package cli

import (
	"fmt"
	"os"

	"github.com/raibru/gsnet/cmd/pktservice/etc"
	"github.com/raibru/gsnet/internal/archive"
	"github.com/raibru/gsnet/internal/pkt"
	"github.com/raibru/gsnet/internal/service"
	"github.com/raibru/gsnet/internal/sys"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pktservice",
	Short: "Packet Service",
	Long:  `Provide packet transform and transfer service behavior via tcp/ip communication`,
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

	wait := make(chan bool)
	sys.StartSignalHandler()

	if configFile != "" {
		var cf = &etc.PktServiceConfig{}
		err := cf.LoadConfig(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error read config file %s: %s\n", configFile, err.Error())
			sys.Exit(2)
		}

		lp := &sys.LoggingParam{
			Service:    cf.Service.Name,
			Version:    VersionShort(),
			Filename:   cf.Logging.Filename,
			TimeFormat: cf.Logging.TimeFormat,
			TeeStdout:  teeStdout,
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

		var arch *archive.Archive
		var archivate chan *archive.Record

		if cf.Archive.Use {
			archivate = make(chan *archive.Record, 10)
			arch = archive.NewArchive(cf.Archive.Filename, cf.Archive.Type, cf.Service.Name)
			arch.SetArchivate(archivate)
			arch.Start(wait)
		}

		for _, elem := range cf.Service.Network {

			cliService := service.NewClientService(
				elem.Channel.Dialer.Name,
				elem.Channel.Dialer.Host,
				elem.Channel.Dialer.Port,
				elem.Channel.Dialer.Retry)

			srvService := service.NewServerService(
				elem.Channel.Listener.Name,
				elem.Channel.Listener.Host,
				elem.Channel.Listener.Port)

			transfer := make(chan []byte)
			srvService.SetForward(transfer)
			cliService.SetTransfer(transfer)

			notify := make(chan []byte)
			cliService.SetReceive(notify)
			srvService.SetProcess(notify)

			pktService := service.NewPacketService(
				elem.Channel.Name,
				elem.Channel.Type)

			pktService.SetDialer(cliService)
			pktService.SetListener(srvService)

			if cf.Archive.Use {
				cliService.SetArchivate(archivate)
				srvService.SetArchivate(archivate)
				pktService.SetArchivate(archivate)
			}

			go func() {
				err := pktService.ApplyConnection()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Fatal Failure. See log. Exit service: %s\n", err.Error())
					sys.Exit(2)
				}
				<-pktService.Mode
			}()
		}
	}

	<-wait

	return nil
}
