package app

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type config struct {
	Port int
	Env  string
}

type Application struct {
	Config config
	Logger *slog.Logger
}

func NewApplication() (*Application, error) {
	var cfg config
	flag.IntVar(&cfg.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Env, "env", "dev", "Environment (dev|prod)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &Application{
		Config: cfg,
		Logger: logger,
	}
	return app, nil
}

func (app *Application) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "available",
		"env":    app.Config.Env,
	})
}
