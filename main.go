package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/weedbox/common-modules/configs"
	"github.com/weedbox/common-modules/daemon"
	"github.com/weedbox/common-modules/http_server"
	"github.com/weedbox/common-modules/logger"
	"github.com/weedbox/common-modules/nats_connector"
	"github.com/weedbox/service-boilerplate/pkg/apis"

	"go.uber.org/fx"
)

var config *configs.Config
var host string
var port int

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
		http_server.Module("web_service"),
		apis.Module("customized_apis"),
		daemon.Module("daemon"),
		fx.NopLogger,
	)

	app.Run()

	return nil
}
