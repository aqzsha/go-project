package hall

import (
	"movies/test/app/core/header"
	api "movies/test/app/core/request"
	"movies/test/factories"
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
)

func TestHall_Create_Success(t *testing.T) {
	req := api.NewRequest(t)

	authUser := header.AuthUser(req.Server.DB())

	cinema := factories.CreateCinema(
		req.Server.DB(),
		factories.CreateCinemaDetails(req.Server.DB()),
	)

	body := map[string]any{
		"cinema_id": cinema.ID,
		"name":      faker.Word(),
		"seats":     int64(100),
	}

	res := req.Post(
		"/hall/store",
		http.StatusOK,
		&body,
		authUser,
	)

	res.Value("success").Boolean().IsTrue()

	data := res.Value("data").Object()

	data.Value("id").Number().Gt(0)
	data.Value("cinema_id").Number().IsEqual(cinema.ID)
	data.Value("name").String().NotEmpty()
	data.Value("seats").Number().Gt(0)

	res.Value("message").String().IsEqual("OK")
}
