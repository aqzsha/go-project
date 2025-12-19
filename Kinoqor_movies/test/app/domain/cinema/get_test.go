package cinema

import (
	"fmt"
	"movies/test/app/core/header"
	api "movies/test/app/core/request"
	"movies/test/factories"
	"net/http"
	"testing"
)
func TestCinema_Get_Success(t *testing.T) {
	req := api.NewRequest(t)

	authUser := header.AuthUser(req.Server.DB())

	cinema := factories.CreateCinema(req.Server.DB(), factories.CreateCinemaDetails(req.Server.DB()))

	res := req.Get(
		fmt.Sprintf("/cinema/get/%d", cinema.ID),
		http.StatusOK,
		nil,
		authUser,
	)

	res.Value("success").Boolean().IsTrue()

	data := res.Value("data").Object()

	data.Value("id").Number()
	data.Value("name").String()
	data.Value("details_id").Number()

	details := data.Value("details").Object()

	details.Value("id").Number()
	details.Value("name").String()
	details.Value("description").String()
	details.Value("address").String()
	details.Value("latitude").Number()
	details.Value("longitude").Number()
	details.Value("created_at").String()

	res.Value("message").String().IsEqual("OK")
}
