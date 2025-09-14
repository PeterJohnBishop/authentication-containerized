package server

import (
	"log"
	"redesigned-telegram/db"
	"redesigned-telegram/server/handlers"
	"redesigned-telegram/server/middleware"

	"github.com/gin-gonic/gin"
)

func StartAuthenticationServer() {
	// Connect to PostgreSQL
	postgres := db.ConnectPSQL()
	err := postgres.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer postgres.Close()

	handlers.CreateUsersTable(postgres)

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	protected := router.Group("/api")
	protected.Use(middleware.JWTMiddleware())
	addOpenRoutes(router, postgres)
	addProtectedRoutes(protected, postgres)
	log.Println("[CONNECTED] authentication server on :8081")
	router.Run(":8081")
}
