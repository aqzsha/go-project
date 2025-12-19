package booking

type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type StoreSeat struct {
	HallID int64 `json:"hall_id"`
	Row    int64 `json:"row"`
	Number int64 `json:"number"`
}
