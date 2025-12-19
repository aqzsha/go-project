package requests

type LoginDTO struct {
	Email     string
	Password  string
	FACode    string
	Lang      string
	UserIp    string
	UserAgent string
}

type RefreshTokenDTO struct {
	RefreshToken string
}
