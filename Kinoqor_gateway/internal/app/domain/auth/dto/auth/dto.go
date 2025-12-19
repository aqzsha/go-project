package auth

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthUser struct {
	ID         int64  `json:"id"`          //X-User-Id
}

