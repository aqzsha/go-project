package user

type CreateDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type DeleteDTO struct {
	Email string `json:"email" validate:"required,email"`
}
