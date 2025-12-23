package server

import (
	"log"
	"net/http"
	"yangdongju/gtd_todo/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ginAdapter struct {
	userHandler *user.UserHandler
}

func (a *ginAdapter) signUp(c *gin.Context) {
	handleJSONRequest(c, &user.SignUpRequest{}, a.userHandler.HandleSignUp)
}

func (a *ginAdapter) login(c *gin.Context) {
	handleJSONRequest(c, &user.LoginRequest{}, a.userHandler.HandleLogin)
}

func handleJSONRequest[T any, R any](c *gin.Context, payload *T, handle func(T) (int, R)) {
	if err := c.ShouldBindJSON(payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	code, res := handle(*payload)
	c.JSON(code, res)
	log.Printf("API response : status=%v / path=%v / res=%v", code, c.Request.RequestURI, res)
}

func SetupRouter(pool *sqlx.DB) *gin.Engine {
	router := gin.Default()
	ginAdapter := ginAdapter{
		userHandler: user.IntializeHandler(pool),
	}
	router.GET("/api/health", healthHandler)
	router.POST("/api/auth/signup", ginAdapter.signUp)
	router.POST("/api/auth/login", ginAdapter.login)

	return router
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}