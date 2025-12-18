package film

import (
	"movies/test/app/core/header"
	api "movies/test/app/core/request"
	"net/http"
	"testing"
	"time"

	"math/rand"

	"github.com/go-faker/faker/v4"
)

func TestFilm_Create_Success(t *testing.T) {
	req := api.NewRequest(t)

	RandNum, _ := faker.RandomInt(1, 10)
	authUser := header.AuthUser(req.Server.DB())

	body := map[string]any{
		"age_limit":  []int{0, 6, 12, 16, 18}[rand.Intn(5)],
		"description": faker.Paragraph(),
		"details_id":  RandNum[0],
		"director":   faker.Name(),
		"duration":   "2h 30m",
		"end_date":   time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
		"name":       faker.Word(),
		"premier":    time.Now().AddDate(0, 0, rand.Intn(30)).Format("2006-01-02"),
		"production": faker.Word(),
		"rate":       rand.Float64()*9 + 1, 
		"start_date": time.Now().Format("2006-01-02"),
	}



	res := req.Post(
		"/film/store",
		http.StatusOK,
		&body,
		authUser,
	)

	res.Value("success").Boolean().IsTrue()

	data := res.Value("data").Object()

	data.Value("id").Number().Gt(0)
	data.Value("name").String().NotEmpty()
	data.Value("description").String()
	data.Value("details_id").Number().Gt(0)

	data.Value("start_date").String().NotEmpty()
	data.Value("end_date").String().NotEmpty()
	data.Value("created_at").String().NotEmpty()

	details := data.Value("details").Object()

	details.Value("id").Number().IsEqual(data.Value("details_id").Number().Raw())
	details.Value("age_limit").Number()
	details.Value("director").String().NotEmpty()
	details.Value("duration").String().NotEmpty()
	details.Value("premier").String().NotEmpty()
	details.Value("production").String().NotEmpty()
	details.Value("rate").Number()

	res.Value("message").String().IsEqual("OK")
}

func TestFilm_Create_ValidationError(t *testing.T) {
	req := api.NewRequest(t)

	RandNum, _ := faker.RandomInt(1, 10)
	authUser := header.AuthUser(req.Server.DB())

	body := map[string]any{
		"age_limit":  []int{0, 6, 12, 16, 18}[rand.Intn(5)],
		"description": faker.Paragraph(),
		"details_id":  RandNum[0],
		"director":   faker.Name(),
		"duration":   "2h 30m",
		"end_date":   time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
		// "name":       faker.Word(),
		"premier":    time.Now().AddDate(0, 0, rand.Intn(30)).Format("2006-01-02"),
		"production": faker.Word(),
		"rate":       rand.Float64()*9 + 1, 
		"start_date": time.Now().Format("2006-01-02"),
	}


	res := req.Post(
		"/film/store",
		http.StatusUnprocessableEntity,
		&body,
		authUser,
	)

	res.Value("success").Boolean().IsFalse()
	res.Value("data").IsNull()

	message := res.Value("message").Object()
	message.Value("name").String().IsEqual("field is required")
}