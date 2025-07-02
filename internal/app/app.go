package app

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/y3933y3933/joker/internal/api"
	"github.com/y3933y3933/joker/internal/db/sqlc"
	"github.com/y3933y3933/joker/internal/service"
	"github.com/y3933y3933/joker/internal/store"
	"github.com/y3933y3933/joker/internal/ws"
)

type config struct {
	Port   int
	Env    string
	DB_URL string
}

type db struct {
	ConnPool *pgxpool.Pool
	Queries  *sqlc.Queries
}

type Application struct {
	Config        config
	Logger        *slog.Logger
	DB            *db
	GameHandler   *api.GameHandler
	GameStore     *store.PostgresGameStore
	PlayerHandler *api.PlayerHandler
	RoundHandler  *api.RoundHandler
	WSHandler     *ws.Handler
}

func NewApplication() (*Application, error) {
	var cfg config
	flag.IntVar(&cfg.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Env, "env", "dev", "Environment (dev|prod)")
	flag.StringVar(&cfg.DB_URL, "db", "", "database url")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	pgDB, queries, err := store.Open(cfg.DB_URL)
	if err != nil {
		return nil, err
	}

	// store
	gameStore := store.NewPostgresGameStore(queries)
	playerStore := store.NewPostgresPlayerStore(queries)
	roundStore := store.NewPostgresRoundStore(queries)
	questionStore := store.NewPostgresQuestionStore(queries)

	// service
	gameService := service.NewGameService(gameStore, playerStore)
	playerService := service.NewPlayerService(playerStore, gameStore)
	roundService := service.NewRoundService(roundStore, playerStore, gameStore)
	questionService := service.NewQuestionService(questionStore)

	// ws
	hub := ws.NewHub()

	// handler
	gameHandler := api.NewGameHandler(gameService, questionService, hub, logger)
	playerHandler := api.NewPlayerHandler(playerService, hub, logger)
	roundHandler := api.NewRoundHandler(roundService, logger, hub)
	wsHandler := ws.NewHandler(hub, logger, playerService, gameService, roundService)

	app := &Application{
		Config: cfg,
		Logger: logger,
		DB: &db{
			ConnPool: pgDB,
			Queries:  queries,
		},
		GameHandler:   gameHandler,
		GameStore:     gameStore,
		PlayerHandler: playerHandler,
		RoundHandler:  roundHandler,
		WSHandler:     wsHandler,
	}
	return app, nil
}

func (app *Application) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "available",
		"env":    app.Config.Env,
	})
}
