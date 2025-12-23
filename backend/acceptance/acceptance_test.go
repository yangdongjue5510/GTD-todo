package acceptance

import (
	"net/http/httptest"
	"testing"

    "yangdongju/gtd_todo/internal/server"
	"yangdongju/gtd_todo/testhelper"
)

func TestUserSenarioes(t *testing.T) {
	testhelper.CleanUp()

	router := server.SetupRouter(testhelper.GetTestDB())
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	apiDriver := apiDriverImpl{
		client:  server.Client(),
		baseURL: server.URL,
	}

	userDsl := userDslImpl{apiDriver: apiDriver}
	SignUpSenario(t, userDsl)
}
