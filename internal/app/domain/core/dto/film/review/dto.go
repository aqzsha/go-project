package review

type CreateReviewDTO struct {
	FilmID int64   `json:"film_id" validate:"required"`
	Body   string  `json:"body"`
	Rating float64 `json:"rating"`
	Title string 	`json:"title" validate:"required"`
}

type ServiceCreateReviewDTO struct {
	FilmID int64  
	UserID int64   
	Title string
	Body   string  
	Rating float64 
}
