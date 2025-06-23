package app

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/y3933y3933/joker/internal/api"
	"github.com/y3933y3933/joker/internal/database"
	"github.com/y3933y3933/joker/internal/ws"
)

type config struct {
	Port   int
	Env    string
	DB_URL string
}

type Application struct {
	Logger         *slog.Logger
	DBQueries      *database.Queries
	DB             *pgxpool.Pool
	Config         config
	GamesHandler   *api.GamesHandler
	PlayersHandler *api.PlayersHandler
	RoundsHandler  *api.RoundsHandler
	WSHub          *ws.Hub
}

func NewApplication() (*Application, error) {
	var cfg config

	flag.IntVar(&cfg.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Env, "env", "dev", "Environment (dev|prod)")
	flag.StringVar(&cfg.DB_URL, "db-url", "", "DATABASE URL")

	flag.Parse()

	loggerHandler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(loggerHandler)

	dbpool, err := pgxpool.New(context.Background(), cfg.DB_URL)
	if err != nil {
		logger.Error("Unable to create connection pool:", err)
		os.Exit(1)
	}
	err = dbpool.Ping(context.Background())
	if err != nil {
		logger.Error("Unable to create connection pool:", err)
		os.Exit(1)
	}

	queries := database.New(dbpool)

	hub := ws.NewHub()
	go hub.Run()

	// handler
	gamesHandler := api.NewGamesHandler(queries, logger, hub)
	playersHandler := api.NewPlayersHandler(queries, logger, hub)
	roundsHandler := api.NewRoundsHandler(queries, logger, hub)

	app := &Application{
		Logger:         logger,
		DB:             dbpool,
		DBQueries:      queries,
		Config:         cfg,
		GamesHandler:   gamesHandler,
		PlayersHandler: playersHandler,
		RoundsHandler:  roundsHandler,
		WSHub:          hub,
	}

	return app, nil
}

func (app *Application) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "available",
		"env":     "dev",
		"version": "1.0.0",
	})
}

func (app *Application) Close() {
	app.DB.Close()
}
