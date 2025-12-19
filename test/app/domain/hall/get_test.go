package hall

import (
	"fmt"
	"movies/test/app/core/header"
	api "movies/test/app/core/request"
	"movies/test/factories"
	"net/http"
	"testing"
)

func TestHall_Get_Success(t *testing.T) {
	req := api.NewRequest(t)

	authUser := header.AuthUser(req.Server.DB())

	// зависимости
	cinema := factories.CreateCinema(req.Server.DB(), factories.CreateCinemaDetails(req.Server.DB()))
	hall := factories.CreateHall(req.Server.DB(), cinema.ID)

	res := req.Get(
		fmt.Sprintf("/hall/get/%d", hall.ID),
		http.StatusOK,
		nil,
		authUser,
	)

	res.Value("success").Boolean().IsTrue()

	data := res.Value("data").Object()

	data.Value("id").Number()
	data.Value("cinema_id").Number().IsEqual(cinema.ID)
	data.Value("name").String()
	data.Value("seats").Number()

	res.Value("message").String().IsEqual("OK")
}
