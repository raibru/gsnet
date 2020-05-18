package cli

import (
	"fmt"
	"os"

	"github.com/raibru/gsnet/cmd/pktservice/etc"
	"github.com/raibru/gsnet/internal/arch"
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
			Version:    Version,
			Filename:   cf.Logging.Filename,
			TimeFormat: cf.Logging.TimeFormat,
		}

		if err := sys.InitLogging(lp); err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error initialize logging: %s\n", err.Error())
			sys.Exit(2)
		}

		loggables := []sys.LoggableContext{
			sys.LogContext, service.LogContext, pkt.LogContext, arch.LogContext,
		}

		for _, c := range loggables {
			if err := c.ApplyLogger(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: apply logging for: %s -> %s\n", c.GetContextName(), err.Error())
			}
		}

		archive := arch.NewArchive(cf.Archive.Filename, cf.Archive.Type, cf.Service.Name)
		archive.Start()

		for _, elem := range cf.Service.Network {

			fmt.Fprintf(os.Stdout, "###==> Info: iterate elem: %s\n", elem.Channel.Name)
			var pktService = new(service.PacketServiceData)
			pktService.Name = elem.Channel.Name
			pktService.Type = elem.Channel.Type
			pktService.Archive = archive
			pktService.Mode = make(chan string)

			var cliService service.ClientServiceData
			cliService.Name = elem.Channel.Dialer.Name
			cliService.Addr = elem.Channel.Dialer.Host
			cliService.Port = elem.Channel.Dialer.Port
			cliService.Transfer = make(chan []byte)
			cliService.Arch = archive
			pktService.Dialer = cliService

			var srvService service.ServerServiceData
			srvService.Name = elem.Channel.Listener.Name
			srvService.Addr = elem.Channel.Listener.Host
			srvService.Port = elem.Channel.Listener.Port
			srvService.Transfer = cliService.Transfer
			srvService.Arch = archive
			pktService.Listener = srvService

			go func() {
				err := pktService.ApplyConnection()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Fatal Failure. See log. Exit service: %s\n", err.Error())
					sys.Exit(2)
				}
				<-pktService.Mode
			}()
		}
		wait := make(chan string)
		<-wait
	}

	return nil
}
