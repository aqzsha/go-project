package film

type CreateHallDTO struct {
	CinemaID int64  `json:"cinema_id"`
	Name     string `json:"name"  validate:"required"`
	Seats    int64  `json:"seats"  validate:"required"`
}
