package cmd

import (
	exServer "github.com/anhdt-vnpay/f5_dynamic_gateway/examples/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var RootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  "",
}

var StartPingExample = &cobra.Command{
	Use:   "example_ping_start_api",
	Short: "Start Ping Example API",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		exServer.GatewayGrpcServer(port)
	},
}

var StartPingExampleGateway = &cobra.Command{
	Use:   "example_ping_start_gateway",
	Short: " Start ewallet service",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		exServer.GatewayServer(port)
	},
}

func init() {
	StartPingExample.Flags().Int("port", 9000, "port")
	StartPingExampleGateway.Flags().Int("port", 8900, "port")
	//add command
	RootCmd.AddCommand(StartPingExample)

	RootCmd.AddCommand(StartPingExampleGateway)
}
