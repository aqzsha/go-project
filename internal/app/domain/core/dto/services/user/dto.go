package user

type CreateDTO struct {
	Email    string
	Password string
}

type DeleteDTO struct {
	Email string
}
