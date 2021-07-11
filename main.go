package main

import (
	"os"

	"github.com/iwanhae/kubenchctl/pkg/tools/network_client_bandwidth"
	"github.com/iwanhae/kubenchctl/pkg/tools/network_client_tps"
	"github.com/iwanhae/kubenchctl/pkg/tools/network_server"
	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use: "kubenchctl",
	}
	cmd.AddCommand(toolCMD())
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func toolCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tool",
		Aliases: []string{"tools", "t"},
	}
	cmd.AddCommand(network_server.NetworkServerCMD())
	cmd.AddCommand(network_client_bandwidth.NetworkClientCMD())
	cmd.AddCommand(network_client_tps.NetworkClientCMD())
	return cmd
}
