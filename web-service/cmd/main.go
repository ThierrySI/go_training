package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_training/web-service/pkg/album"
	"go_training/web-service/pkg/config"
	"go_training/web-service/pkg/database"
)

func main() {

	// 1. Read Config
	config := config.LoadConfigFile()

	// 2. DB Pool
	db, _ := database.ConnectCH(config.DBHostname, config.DBPort, config.DBUsername, config.DBPassword, config.DBDatabase, config.DBMaxOpenConn, config.DBMaxIdleConn)
	defer database.CloseDB(db)

	// 3. Create Gin Router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 4. Register routes
	album.RegisterRoutes(router, db)

	// 5. Attached the router to httpServer and start it on :8080
	router.Run(fmt.Sprintf("%v:%v", config.HTTPServer, config.HTTPPort))

}
