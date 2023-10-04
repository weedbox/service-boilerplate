package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/weedbox/common-modules/configs"
	"github.com/weedbox/common-modules/daemon"
	"github.com/weedbox/common-modules/healthcheck_apis"
	"github.com/weedbox/common-modules/http_server"
	"github.com/weedbox/common-modules/logger"
	"github.com/weedbox/common-modules/nats_connector"
	"github.com/weedbox/common-modules/postgres_connector"
	"github.com/weedbox/service-boilerplate/pkg/apis"
	"github.com/weedbox/service-boilerplate/pkg/apis_db"
	"github.com/weedbox/service-boilerplate/pkg/example_rpc"
	"github.com/weedbox/websocket-modules/system_rpc"
	"github.com/weedbox/websocket-modules/websocket_endpoint"
	"github.com/weedbox/websocket-modules/websocket_server"

	"go.uber.org/fx"
)

var config *configs.Config
var host string
var port int
var pc bool

var rootCmd = &cobra.Command{
	Use:   "service-boilerplate",
	Short: "general service",
	Long:  `service-boilerplate is a general service.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := run(); err != nil {
			return err
		}
		return nil
	},
}

func init() {

	config = configs.NewConfig("SERVICE")

	rootCmd.Flags().StringVar(&host, "host", "0.0.0.0", "Specify host")
	rootCmd.Flags().IntVar(&port, "port", 8080, "Specify service port")
	rootCmd.Flags().BoolVar(&pc, "print_configs", false, "Show all configs")
}

func main() {

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run() error {

	config.SetConfigs(map[string]interface{}{
		"web_service.host": host,
		"web_service.port": port,
	})

	app := fx.New(
		fx.Supply(config),

		// Modules
		logger.Module(),
		nats_connector.Module("internal_event"),
		postgres_connector.Module("database"),
		http_server.Module("web_service"),
		healthcheck_apis.Module("healthcheck_apis"),

		// Websocket
		websocket_server.Module("websocket_server"),
		websocket_endpoint.Module("endpoint_example", "/example"),
		system_rpc.Module("rpc_system", "/example"),

		// Customization
		apis.Module("customized_apis"),
		apis_db.Module("customized_apis_db"),
		example_rpc.Module("example_rpc"),

		// Integration
		daemon.Module("daemon"),
		fx.NopLogger,
	)

	if pc {
		config.PrintAllSettings()
	}

	app.Run()

	return nil
}
