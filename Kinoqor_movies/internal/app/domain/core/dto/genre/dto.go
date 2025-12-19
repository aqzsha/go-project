package genre

type CreateGenreDTO struct {
	Name        string    `json:"name"  validate:"required"`
}
