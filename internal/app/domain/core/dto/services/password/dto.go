package password

type ForgotPasswordDTO struct {
	Email string
}

type ResetPasswordDTO struct {
	Token       string
	NewPassword string
}

type ChangePasswordDTO struct {
	NewPassword     string
	Password        string
	ConfirmPassword string
	Token           string
}
