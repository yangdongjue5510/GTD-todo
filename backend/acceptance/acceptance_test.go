package acceptance

import (
	"net/http/httptest"
	"testing"

	"yangdongju/gtd_todo/internal/server"
	"yangdongju/gtd_todo/testhelper"

	"github.com/gin-gonic/gin"
)

func TestUserSenarioes(t *testing.T) {
	router := server.SetupRouter(testhelper.GetTestDB())
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	gin.SetMode(gin.TestMode)
	
	apiDriver := apiDriverImpl{
		client:  server.Client(),
		baseURL: server.URL,
	}

	userDsl := userDslImpl{apiDriver: apiDriver}
	signUpSenario(t, userDsl)
	loginSenario(t, userDsl)
}
