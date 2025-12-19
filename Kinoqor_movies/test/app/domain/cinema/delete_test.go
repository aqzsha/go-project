package cinema

import (
	"fmt"
	"movies/test/app/core/header"
	api "movies/test/app/core/request"
	"movies/test/factories"
	"net/http"
	"testing"
)

func TestCinema_Delete_Success(t *testing.T) {
	req := api.NewRequest(t)

	authUser := header.AuthUser(req.Server.DB())

	cinema := factories.CreateCinema(req.Server.DB(), factories.CreateCinemaDetails(req.Server.DB()))

	res := req.Delete(
		fmt.Sprintf("/cinema/delete/%d", cinema.ID),
		http.StatusOK,
		nil,
		authUser,
	)

	res.Value("success").Boolean().IsTrue()

	res.Value("data").Boolean().IsTrue()

	res.Value("message").String().IsEqual("deleted")
}
