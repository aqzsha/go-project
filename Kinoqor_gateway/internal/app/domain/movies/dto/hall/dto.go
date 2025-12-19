package hall

type CreateHallDTO struct {
	CinemaID int64  `json:"cinema_id"`
	Name     string `json:"name"`
	Seats    int64  `json:"seats"`
}
