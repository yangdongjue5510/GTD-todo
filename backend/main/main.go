package main

import (
	"log"
	"net/http"
	"yangdongju/gtd_todo/internal/config"
	"yangdongju/gtd_todo/internal/db"
	"yangdongju/gtd_todo/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.Load()
	pool, err := db.NewConnectionPool(cfg)
	if err != nil {
		log.Fatalf("Database connection failed.\n%v", err)
	}
	setupRouter(pool)
}

func setupRouter(pool *sqlx.DB) *gin.Engine {
	router := gin.Default()
	ginAdapter := ginAdapter{
		userHandler: user.IntializeHandler(pool),
	}
	router.GET("/api/health", healthHandler)
	router.POST("/api/auth/signup", ginAdapter.signUp)

	return router
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
