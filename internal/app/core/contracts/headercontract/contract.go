package headercontract

import (
	"context"
)

const (
	UserIdKey         = "X-User-Id"
)

type BearerTokenKey struct{}
type AuthUserKey struct{}

type AuthUser struct {
	ID         int64  //X-User-Id
}

func GetBearer(cxt context.Context) string {
	return cxt.Value(BearerTokenKey{}).(string)
}
