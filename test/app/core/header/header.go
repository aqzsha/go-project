package header

import (
	"movies/test/factories"
	"strconv"

	header "movies/internal/app/core/contracts/microservices/header-contract"

	"github.com/gavv/httpexpect/v2"
	"gorm.io/gorm"
)

func ApplyAuth(req *httpexpect.Request, u header.AuthUser) *httpexpect.Request {
	req = req.
		WithHeader(header.UserIdKey, strconv.FormatInt(u.ID, 10))

	return req
}

func AuthUser(db *gorm.DB) header.AuthUser {
	user := factories.CreateUser(db)

	return header.AuthUser{
		ID:         user.ID,
	}
}