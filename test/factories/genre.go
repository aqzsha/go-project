package factories

import (
	"movies/internal/app/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func CreateGenre(db *gorm.DB) *models.Genre {
	name := faker.Word()

	genre := &models.Genre{
		Name: name,
	}

	db.Create(genre)
	return genre
}
