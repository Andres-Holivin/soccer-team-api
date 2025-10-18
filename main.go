package main

import (
	"log"
	"soccer-team-api/router"
	"soccer-team-api/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	env := utils.LoadEnv()
	gin.SetMode(env.GinMode)
	utils.InitCloudinary()

	if err := utils.ConnectDB(env); err != nil {
		log.Fatalf("❌ Failed to connect database: %v", err)
	}
	// if err := utils.MigrateDB(); err != nil {
	// 	log.Fatalf("❌ Failed to migrate database: %v", err)
	// }
	r := router.SetupRouter()
	r.SetTrustedProxies([]string{env.AllowOrigins})
	log.Println("✅ Server started on port:", env.AppPort)
	if err := r.Run(":" + env.AppPort); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
