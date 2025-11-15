package password

type ForgotPasswordDTO struct {
	Email string
}

type ResetPasswordDTO struct {
	Email       string
	Token       string
	NewPassword string
}

type ChangePasswordDTO struct {
	NewPassword     string
	Password        string
	ConfirmPassword string
	Token           string
}

