package main

import (
	"log"
	"net/http"
	"yangdongju/gtd_todo/internal/user"

	"github.com/gin-gonic/gin"
)

type ginAdapter struct {
	userHandler *user.UserHandler
}

func (a *ginAdapter) signUp(c *gin.Context) {
	var req user.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	
	code, res := a.userHandler.HandleSignUp(req)
	c.JSON(code, res)
	log.Printf("API response : status=%v / path=%v / res=%v", code, c.Request.RequestURI, res)

}