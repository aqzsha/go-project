package cinema

import (
	"movies/test/app/core/header"
	api "movies/test/app/core/request"
	"net/http"
	"testing"

	"github.com/go-faker/faker/v4"
)
func TestCinema_Create_Success(t *testing.T) {
	req := api.NewRequest(t)

	authUser := header.AuthUser(req.Server.DB())

	body := map[string]any{
		"name":        faker.Word(),
		"description": faker.Sentence(),
		"address":     faker.Sentence(),
		"latitude":    51.1694,
		"longitude":   71.4491,
	}

	res := req.Post(
		"/cinema/store",
		http.StatusOK,
		&body,
		authUser,
	)

	res.Value("success").Boolean().IsTrue()

	data := res.Value("data").Object()

	// cinema
	data.Value("id").Number().Gt(0)
	data.Value("name").String()
	data.Value("details_id").Number().Gt(0)

	// cinema details
	details := data.Value("details").Object()

	details.Value("id").Number().IsEqual(
		data.Value("details_id").Number().Raw(),
	)
	details.Value("name").String().NotEmpty()
	details.Value("description").String()
	details.Value("address").String()
	details.Value("latitude").Number()
	details.Value("longitude").Number()
	details.Value("created_at").String().NotEmpty()

	res.Value("message").String().IsEqual("OK")
}
