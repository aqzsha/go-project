package genre

import (
	"fmt"
	"movies/test/app/core/header"
	api "movies/test/app/core/request"
	"movies/test/factories"
	"net/http"
	"testing"
)

func TestGenre_Delete_Success(t *testing.T) {
	req := api.NewRequest(t)

	authUser := header.AuthUser(req.Server.DB())

	genre := factories.CreateGenre(req.Server.DB())

	res := req.Delete(
		fmt.Sprintf("/genre/delete/%d", genre.ID),
		http.StatusOK,
		nil,
		authUser,
	)

	res.Value("success").Boolean().IsTrue()

	res.Value("data").Boolean().IsTrue()

	res.Value("message").String().IsEqual("deleted")
}
