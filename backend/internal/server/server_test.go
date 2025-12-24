package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"yangdongju/gtd_todo/internal/server"
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
	router := server.SetupRouter(testhelper.GetTestDB())

	// when
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/health", nil)
	router.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `{"status":"ok"}`)
}


