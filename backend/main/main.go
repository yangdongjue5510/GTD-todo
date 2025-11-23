package main

import (
	"log"
	"net/http"
	"time"
	"yangdongju/gtd_todo/internal/config"
	"yangdongju/gtd_todo/internal/db"
	"yangdongju/gtd_todo/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	userRepository user.UserRepository
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
	userRepository = user.NewUserRepository(pool)
	router.GET("/api/health", healthHandler)
	router.POST("/api/signup", signupHandler)

	return router
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func signupHandler(c *gin.Context) {
	var req SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Sign up API failed: %v\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foundUser, err := userRepository.FindUserByEmail(req.Email)
	if err != nil {
		log.Printf("Sign up API failed : %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sign up failed"})
		return
	}

	if foundUser != nil {
		log.Printf("Sign up API failed: User exists with inputed email.%v", foundUser)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User exists with inputed email."})
		return
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		log.Printf("Sign up API failed : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sign up failed"})
		return
	}

	newUser := user.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	savedUser, err := userRepository.Save(&newUser)
	if err != nil {
		log.Printf("Sign up API failed : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sign up failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    savedUser.ID,
		"email": savedUser.Email,
	})
}

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
