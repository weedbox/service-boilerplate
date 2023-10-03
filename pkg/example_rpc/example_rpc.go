package example_rpc

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/weedbox/common-modules/http_server"
	"github.com/weedbox/websocket-modules/websocket_server"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ExampleRPC struct {
	params Params
	logger *zap.Logger
	router *gin.RouterGroup
	scope  string
}

type Params struct {
	fx.In

	Lifecycle       fx.Lifecycle
	Logger          *zap.Logger
	HTTPServer      *http_server.HTTPServer
	WebSocketServer *websocket_server.WebSocketServer
}

func Module(scope string) fx.Option {

	var erpc *ExampleRPC

	return fx.Module(
		scope,
		fx.Provide(func(p Params) *ExampleRPC {

			erpc := &ExampleRPC{
				params: p,
				logger: p.Logger.Named(scope),
				scope:  scope,
			}

			return erpc
		}),
		fx.Populate(&erpc),
		fx.Invoke(func(p Params) {

			p.Lifecycle.Append(
				fx.Hook{
					OnStart: erpc.onStart,
					OnStop:  erpc.onStop,
				},
			)
		}),
	)
}

func (erpc *ExampleRPC) onStart(ctx context.Context) error {

	erpc.logger.Info("Starting Example RPC")

	ep := erpc.params.WebSocketServer.GetEndpoint("/example")
	if ep == nil {
		return errors.New("Not found endpoint")
	}

	ep.GetAdapter().Register("Example.Hello", erpc.exampleHello)

	return nil
}

func (erpc *ExampleRPC) onStop(ctx context.Context) error {

	erpc.logger.Info("Stopped Example RPC")

	return nil
}

type Hello struct {
	Name string                 `json:"name"`
	Map  map[string]interface{} `json:"map"`
}

func (erpc *ExampleRPC) exampleHello(c *websocket_server.Context) (interface{}, error) {

	parameters := c.GetRequest().Params.([]interface{})

	data := &Hello{
		Name: parameters[0].(string),
		Map: map[string]interface{}{
			"key1": "Value1",
			"key2": "Value2",
		},
	}

	return data, nil
}
