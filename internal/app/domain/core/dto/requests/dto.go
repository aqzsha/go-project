package requests

type LoginDTO struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,password"`
	FACode    string `json:"2fa_code"`
	Lang      string `json:"lang"`
	UserIp    string `json:"userIp"`
	UserAgent string `json:"userAgent"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
