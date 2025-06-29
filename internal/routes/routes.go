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
	// 建立遊戲
	games.POST("/", app.GameHandler.HandleCreateGame)

	codes := games.Group("/:code", middleware.ValidateGameExists(app.GameStore))
	{
		// 加入遊戲
		codes.POST("/join", app.PlayerHandler.HandleJoinGame)
		// 查看所有玩家
		codes.GET("/players", app.PlayerHandler.HandleListPlayers)
		// 開始遊戲
		codes.POST("/start", app.RoundHandler.HandleStartGame)

		// 取得隨機題目
		codes.GET("/questions", app.GameHandler.HandleGetQuestions)

		// 更新題目
		codes.POST("/rounds/:id/question", app.RoundHandler.HandleSubmitQuestion)

		// 更新回答
		codes.POST("/rounds/:id/answer", app.RoundHandler.HandleSubmitAnswer)

	}
	// ws
	router.GET("/ws/games/:code", app.WSHandler.ServeWS)

	return router
}
