package genre

import (
	"fmt"
	"movies/test/app/core/header"
	api "movies/test/app/core/request"
	"movies/test/factories"
	"net/http"
	"testing"
)

func TestGenre_Get_Success(t *testing.T) {
	req := api.NewRequest(t)

	authUser := header.AuthUser(req.Server.DB())

	genre := factories.CreateGenre(req.Server.DB())

	res := req.Get(
		fmt.Sprintf("/genre/get/%d", genre.ID),
		http.StatusOK,
		nil,
		authUser,
	)

	res.Value("success").Boolean().IsTrue()

	data := res.Value("data").Object()

	data.Value("id").Number()
	data.Value("name").String()

	res.Value("message").String().IsEqual("OK")
}
