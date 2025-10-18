package router

import (
	"soccer-team-api/controller"
	"soccer-team-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Service is healthy",
		})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controller.RegisterUser)
			auth.POST("/login", controller.LoginUser)
		}

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", controller.GetProfile)

			teams := protected.Group("/teams")
			{
				teams.POST("", controller.CreateTeam)
				teams.GET("", controller.GetTeams)
				teams.GET("/:id", controller.GetTeam)
				teams.PUT("/:id", controller.UpdateTeam)
				teams.DELETE("/:id", controller.DeleteTeam)
			}

			players := protected.Group("/players")
			{
				players.POST("", controller.CreatePlayer)
				players.GET("", controller.GetPlayers)
				players.GET("/:id", controller.GetPlayer)
				players.PUT("/:id", controller.UpdatePlayer)
				players.DELETE("/:id", controller.DeletePlayer)
			}

			matches := protected.Group("/matches")
			{
				matches.POST("", controller.CreateMatch)
				matches.GET("", controller.GetMatches)
				matches.GET("/:id", controller.GetMatch)
				matches.PUT("/:id", controller.UpdateMatch)
				matches.DELETE("/:id", controller.DeleteMatch)
				matches.POST("/:id/result", controller.RecordMatchResult)
			}

			reports := protected.Group("/reports")
			{
				reports.GET("/matches", controller.GetAllMatchReports)
				reports.GET("/matches/:id", controller.GetMatchReport)
			}
		}
	}

	return r
}
