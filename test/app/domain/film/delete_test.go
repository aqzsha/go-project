package film

import (
	"fmt"
	"movies/test/app/core/header"
	api "movies/test/app/core/request"
	"movies/test/factories"
	"net/http"
	"testing"
)

func TestFilm_Delete_Success(t *testing.T) {
	req := api.NewRequest(t)

	authUser := header.AuthUser(req.Server.DB())

	filmDetails := factories.CreateFilmDetails(req.Server.DB())
	film := factories.CreateFilm(req.Server.DB(), filmDetails)

	res := req.Delete(
		fmt.Sprintf("/film/delete/%d", film.ID),
		http.StatusOK,
		nil,
		authUser,
	)

	res.Value("success").Boolean().IsTrue()

 	res.Value("data").Boolean().IsTrue()

	res.Value("message").String().IsEqual("deleted")
}
