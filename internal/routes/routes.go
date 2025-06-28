package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/app"
)

func SetupRoutes(app *app.Application) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/healthz", app.HealthCheck)
	}

	games := api.Group("/games")
	{
		games.POST("/", app.GameHandler.HandleCreateGame)
		games.POST("/:code/join", app.PlayerHandler.HandleJoinGame)

	}

	// ws
	router.GET("/ws/games/:code", app.WSHandler.ServeWS)

	return router
}
