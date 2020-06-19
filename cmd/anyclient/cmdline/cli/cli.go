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
	var clientService *service.ClientService
	archiveService := archive.NonArchive()
	readerService := pkt.NonPacketReader()

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

		clientService = service.NewClientService(cf.Service.Name, cf.Service.Host, cf.Service.Port)

		var archivate chan *archive.Record
		var process chan []byte

		if cf.Archive.Use {
			archivate = make(chan *archive.Record, 10)
			archiveService = archive.NewArchive(cf.Archive.Filename, cf.Archive.Type, cf.Service.Name)
			archiveService.SetArchivate(archivate)
			clientService.SetArchivate(archivate)
			//clientService.SetArchive(archiveService.Archivate)
		}

		if cf.Packet.Use {
			process = make(chan []byte)
			readerService = pkt.NewPacketReader(cf.Packet.Filename, cf.Packet.Wait, waitTransfer)
			readerService.SetSupply(process)
			clientService.SetProcess(process)
			//clientService.SetProcess(readerService.Supply)
		}

	} else {
		clientService = service.NewClientService("anyclient", "129.0.0.1", "30100")
	}

	err := clientService.ApplyConnection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Failure. See log. Exit service: %s\n", err.Error())
		sys.Exit(2)
	}

	defer clientService.Finalize()

	wait := make(chan bool, 1)
	readed := make(chan bool, 1)
	done := make(chan bool, 1)

	if archiveService.Use {
		archiveService.Start(wait)
	}
	if readerService.Use {
		readerService.Start(readed)
		clientService.TransferPackets(done)
		<-done
	} else {
		clientService.ReceivePackets(done)
		<-done
	}

	return nil
}
