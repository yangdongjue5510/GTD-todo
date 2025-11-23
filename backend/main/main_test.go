package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"yangdongju/gtd_todo/testhelper"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	testhelper.TestMain(m)
}

func TestHealthCheck(t *testing.T) {
	// given
	testhelper.CleanUp()
	gin.SetMode(gin.TestMode)
	router := setupRouter(testhelper.GetTestDB())

	// when
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/health", nil)
	router.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `{"status":"ok"}`)
}

func TestSignUpAPI(t *testing.T) {
	// given
	testhelper.CleanUp()
	gin.SetMode(gin.TestMode)
	router := setupRouter(testhelper.GetTestDB())

	requestBody := SignUpRequest{Email: "test@example.com", Password: "password1234"}
	requestJsonBody, _ := json.Marshal(requestBody)

	// when
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(requestJsonBody))
	router.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusCreated, w.Code)
}
