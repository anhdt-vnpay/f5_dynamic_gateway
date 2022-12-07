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

var StartPingEchoGateway = &cobra.Command{
	Use:   "example_ping_start_echo_server",
	Short: " Start echo service",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		exServer.GatewayEchoServer(port)
	},
}

var StartPingEchoGateway2 = &cobra.Command{
	Use:   "example_ping_start_echo_server2",
	Short: " Start echo service",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		exServer.GatewayEchoServer2(port)
	},
}

func init() {
	StartPingExample.Flags().Int("port", 9000, "port")
	StartPingExampleGateway.Flags().Int("port", 8900, "port")
	StartPingEchoGateway.Flags().Int("port", 8901, "port")
	StartPingEchoGateway2.Flags().Int("port", 8902, "port")
	//add command
	RootCmd.AddCommand(StartPingExample)

	RootCmd.AddCommand(StartPingExampleGateway)
	RootCmd.AddCommand(StartPingEchoGateway)
	RootCmd.AddCommand(StartPingEchoGateway2)
}
