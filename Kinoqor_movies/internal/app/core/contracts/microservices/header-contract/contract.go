package headercontract

import (
	"context"
	"errors"
)

const (
	UserIdKey         = "X-User-Id"
)

var AuthUserNotFount = errors.New("user data is not specified in the header")

type AuthUserKey struct{}

type AuthUser struct {
	ID         int64
}

func GetAuthUser(ctx context.Context) (AuthUser, error) {
	authUser, ok := ctx.Value(AuthUserKey{}).(AuthUser)
	if ok {
		return authUser, nil
	} else {
		return AuthUser{}, AuthUserNotFount
	}
}
