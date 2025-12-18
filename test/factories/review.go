package factories

import (
	"math/rand"
	"time"

	"movies/internal/app/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func CreateReview(db *gorm.DB, filmID int64, userID int64) *models.Review {
	body := faker.Sentence()
	rating := 1 + rand.Float64()*9 
	now := time.Now()

	review := &models.Review{
		FilmID:    filmID,
		UserID:    userID,
		Body:      body,
		Rating:    rating,
		CreatedAt: &now,
	}

	db.Create(review)
	return review
}
