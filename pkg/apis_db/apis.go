package apis_db

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weedbox/common-modules/http_server"
	"github.com/weedbox/common-modules/postgres_connector"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type APIs struct {
	params Params
	logger *zap.Logger
	router *gin.RouterGroup
	scope  string
}

type Params struct {
	fx.In

	Lifecycle  fx.Lifecycle
	Logger     *zap.Logger
	HTTPServer *http_server.HTTPServer
	Database   *postgres_connector.PostgresConnector
}

func Module(scope string) fx.Option {

	var a *APIs

	return fx.Options(
		fx.Provide(func(p Params) *APIs {

			a := &APIs{
				params: p,
				logger: p.Logger.Named(scope),
				scope:  scope,
			}

			return a
		}),
		fx.Populate(&a),
		fx.Invoke(func(p Params) {

			p.Lifecycle.Append(
				fx.Hook{
					OnStart: a.onStart,
					OnStop:  a.onStop,
				},
			)
		}),
	)

}

func (a *APIs) onStart(ctx context.Context) error {

	a.logger.Info("Starting APIs")

	// Auto migration
	db := a.params.Database.GetDB()
	db.AutoMigrate(&Entry{})

	// Router
	a.router = a.params.HTTPServer.GetRouter().Group("apis")
	a.router.GET("/v1/db/entries", a.list)

	return nil
}

func (a *APIs) onStop(ctx context.Context) error {
	a.logger.Info("Stopped APIs")

	return nil
}

func (a *APIs) list(c *gin.Context) {

	db := a.params.Database.GetDB()

	var entries []Entry
	results := db.Raw("SELECT * FROM entries").Scan(&entries)
	if results.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": results.Error.Error()})
		return
	}

	if len(entries) == 0 {
		entries = []Entry{}
	}

	c.JSON(http.StatusOK, gin.H{
		"entries": entries,
	})
}
