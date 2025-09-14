package server

import (
	"database/sql"
	"redesigned-telegram/server/handlers"

	"github.com/gin-gonic/gin"
)

func addOpenRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, you've reached the authentication server. Please leave a message after the beep.",
		})
	})
	r.POST("auth/login", func(c *gin.Context) {
		handlers.Login(db, c)
	})
	r.POST("auth/register", func(c *gin.Context) {
		handlers.RegisterUser(db, c)
	})
	r.GET("auth/refresh", func(c *gin.Context) {
		handlers.Refresh(c)
	})
}

func addProtectedRoutes(r *gin.RouterGroup, db *sql.DB) {

	r.GET("/users", func(c *gin.Context) {
		handlers.GetUsers(db, c)
	})
	r.GET("/users/:id", func(c *gin.Context) {
		handlers.GetUserByID(db, c)
	})
	r.PUT("/users", func(c *gin.Context) {
		handlers.UpdateUser(db, c)
	})
	r.DELETE("/users/:id", func(c *gin.Context) {
		handlers.DeleteUserByID(db, c)
	})
}
