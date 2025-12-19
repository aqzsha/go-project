package seat

type CreateSeatDTO struct {
	HallID int64 `json:"hall_id"  validate:"required"`
	Row    int64 `json:"row" validate:"required"`
	Number int64 `json:"number"  validate:"required"`
}
