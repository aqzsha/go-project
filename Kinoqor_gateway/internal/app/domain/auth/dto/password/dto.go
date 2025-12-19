package password

type ForgotPasswordDTO struct {
	Email string `json:"email"`
}

type ResetPasswordDTO struct {
	Email                   string `json:"email"`
	Token                   string `json:"token"`
	NewPassword             string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}

type ChangePasswordDTO struct {
	Password        string `json:"password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"new_password_confirmation"`
}

type VerifyPinDTO struct {
	Email   string `json:"email"`
	PinCode string `json:"pin_code"`
}
