package factories

import (
	"movies/internal/app/models"

	"github.com/go-faker/faker/v4"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB) *models.User {
	user := &models.User{
		Email: faker.Email(),
		Password: faker.Password(),
	}

	db.Create(user)
	return user
}
