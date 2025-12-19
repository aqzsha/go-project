package genre

import (
	"movies/test/app/core/header"
	api "movies/test/app/core/request"
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
)

func TestGenre_Create_Success(t *testing.T) {
	req := api.NewRequest(t)

	authUser := header.AuthUser(req.Server.DB())

	body := map[string]any{
		"name": faker.Word(),
	}

	res := req.Post(
		"/genre/store",
		http.StatusOK,
		&body,
		authUser,
	)

	res.Value("success").Boolean().IsTrue()

	data := res.Value("data").Object()

	data.Value("id").Number().Gt(0)
	data.Value("name").String().NotEmpty()

	res.Value("message").String().IsEqual("OK")
}
