package film

import (
	"fmt"
	"movies/test/app/core/header"
	api "movies/test/app/core/request"
	"movies/test/factories"
	"net/http"
	"testing"
)

func TestFilm_Get_Success(t *testing.T) {
	req := api.NewRequest(t)

	authUser := header.AuthUser(req.Server.DB())

	filmDetails := factories.CreateFilmDetails(req.Server.DB())
	film := factories.CreateFilm(req.Server.DB(), filmDetails)

	res := req.Get(
		fmt.Sprintf("/film/get/%d", film.ID),
		http.StatusOK,
		nil,
		authUser,
	)

	res.Value("success").Boolean().IsTrue()

	data := res.Value("data").Object()

	data.Value("id").Number()
	data.Value("name").String()
	data.Value("description").String()
	data.Value("details_id").Number()

	data.Value("start_date").String()
	data.Value("end_date").String()
	data.Value("created_at").String()

	details := data.Value("details").Object()

	details.Value("id").Number()
	details.Value("age_limit").Number()
	details.Value("director").String()
	details.Value("duration").String()
	details.Value("premier").String()
	details.Value("production").String()
	details.Value("rate").Number()

	res.Value("message").String().IsEqual("OK")
}
