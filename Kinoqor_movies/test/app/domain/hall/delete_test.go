package hall

import (
	"fmt"
	"movies/test/app/core/header"
	api "movies/test/app/core/request"
	"movies/test/factories"
	"net/http"
	"testing"
)

func TestHall_Delete_Success(t *testing.T) {
	req := api.NewRequest(t)

	authUser := header.AuthUser(req.Server.DB())

	cinema := factories.CreateCinema(req.Server.DB(), factories.CreateCinemaDetails(req.Server.DB()))
	hall := factories.CreateHall(req.Server.DB(), cinema.ID)

	res := req.Delete(
		fmt.Sprintf("/hall/delete/%d", hall.ID),
		http.StatusOK,
		nil,
		authUser,
	)

	res.Value("success").Boolean().IsTrue()

	res.Value("data").Boolean().IsTrue()

	res.Value("message").String().IsEqual("deleted")
}
