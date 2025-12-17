package film

import "time"

type CreateFilmDTO struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DetailsID   int64     `json:"details_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Duration    string    `json:"duration"`
	Premier     time.Time `json:"premier"`
	Production  string    `json:"production"`
	Director    string    `json:"director"`
	Rate        float32   `json:"rate"`
	AgeLimit    int8      `json:"age_limit"`
}
