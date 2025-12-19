package film

type CreateCinemaDTO struct {
	Name        string  `json:"name" validate:"required,max=255"`
	Description string  `json:"description" validate:"omitempty"`
	Address     string  `json:"address" validate:"required,max=255"`
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
}
