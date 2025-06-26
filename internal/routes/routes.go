package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/joker/internal/app"
	"github.com/y3933y3933/joker/internal/ws"
)

func SetRoutes(app *app.Application) *gin.Engine {
	router := gin.Default()

	router.GET("/api/healthz", app.HealthCheck)

	// games
	games := router.Group("/api/games")
	{

		// 建立遊戲
		games.POST("/", app.GamesHandler.CreateGame)
		// 加入遊戲
		games.POST("/:code/join", app.PlayersHandler.JoinGame)
		// 查看所有玩家
		games.GET("/:code/players", app.PlayersHandler.ListPlayers)

		games.GET("/:code", app.GamesHandler.GetGame)

		// 開始遊戲
		games.POST("/:code/start", app.RoundsHandler.StartGame)

		// 更新題目
		games.POST("/:code/rounds/:id/question", app.RoundsHandler.SubmitQuestion)

		// 更新回答
		games.POST("/:code/rounds/:id/answer", app.RoundsHandler.SubmitAnswer)

		// 抽牌
		games.POST("/:code/rounds/:id/draw", app.RoundsHandler.DrawCard)

		// 換下一輪
		games.POST("/:code/rounds/next", app.RoundsHandler.CreateNextRound)
		// 遊戲結束
		games.POST("/:code/end", app.GamesHandler.EndGame)
		// 踢出玩家
		games.DELETE("/:code/players/:player_id", app.RoundsHandler.RemovePlayer)

		// 選擇題目
		games.GET("/:code/questions", app.GamesHandler.GetQuestions)

	}

	// ws
	router.GET("/ws/games/:code", func(c *gin.Context) {
		ws.ServeWS(app.WSHub, c)
	})

	return router
}
