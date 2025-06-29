package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/app"
	"github.com/y3933y3933/joker/internal/middleware"
)

func SetupRoutes(app *app.Application) *gin.Engine {
	router := gin.Default()

	router.GET("/api/healthz", app.HealthCheck)

	// games
	games := router.Group("/api/games")
	games.POST("/", app.GameHandler.HandleCreateGame)

	codes := games.Group("/:code", middleware.ValidateGameExists(app.GameStore))
	{
		codes.POST("/join", app.PlayerHandler.HandleJoinGame)
		codes.GET("/players", app.PlayerHandler.HandleListPlayers)
		codes.POST("/start", app.RoundHandler.HandleStartGame)
	}
	// ws
	router.GET("/ws/games/:code", app.WSHandler.ServeWS)

	return router
}
