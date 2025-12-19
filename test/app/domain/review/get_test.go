package review

import (
	"fmt"
	"movies/test/app/core/header"
	api "movies/test/app/core/request"
	"movies/test/factories"
	"net/http"
	"testing"
)

func TestFilmReview_Get_Success(t *testing.T) {
	req := api.NewRequest(t)

	authUser := header.AuthUser(req.Server.DB())

	filmDetails := factories.CreateFilmDetails(req.Server.DB())
	film := factories.CreateFilm(req.Server.DB(), filmDetails)
	user := factories.CreateUser(req.Server.DB())

	review := factories.CreateReview(req.Server.DB(), film.ID, user.ID)

	if review.ID == 0 {
		t.Fatalf("review.ID is 0, review was not persisted")
	}

	res := req.Get(
		fmt.Sprintf("/film/review/get/%d", review.ID),
		http.StatusOK,
		nil,
		authUser,
	)

	res.Value("success").Boolean().IsTrue()

	data := res.Value("data").Object()

	data.Value("id").Number()
	data.Value("film_id").Number().IsEqual(film.ID)
	data.Value("user_id").Number().IsEqual(user.ID)
	data.Value("body").String()
	data.Value("rating").Number()

	res.Value("message").String().IsEqual("OK")
}
