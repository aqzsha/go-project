package screening

import (
	dto "booking/internal/app/domain/core/dto/screening"
	"booking/internal/app/models"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("screening not found")

type Repository interface {
	Create(ctx context.Context, input dto.CreateScreeningDTO) (models.Screening, error)
	Get(ctx context.Context, id int64) (models.Screening, error)
	Delete(ctx context.Context, id int64) (bool, error)
	ListByFilm(ctx context.Context, filmId int64) ([]models.Screening, error)
	ListByCinema(ctx context.Context, cinemaId int64) ([]models.Screening, error)
	ListByHall(ctx context.Context, hallId int64) ([]models.Screening, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
func (r *repository) Create(
	ctx context.Context,
	input dto.CreateScreeningDTO,
) (models.Screening, error) {

	var screening models.Screening

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.
			Model(&models.Screening{}).
			Where("hall_id = ?", input.HallID).
			Where("start_at < ? AND end_at > ?", input.EndAt, input.StartAt).
			Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check screening conflicts: %w", err)
		}

		if count > 0 {
			return errors.New("screening time intersects with existing screening in this hall")
		}

		screening = models.Screening{
			HallID:   input.HallID,
			CinemaID: input.CinemaID,
			FilmID:   input.FilmID,
			Price:    input.Price,
			Start_At: input.StartAt,
			End_At:   input.EndAt,
			Language: input.Language,
			Format:   input.Format,
		}

		if err := tx.Create(&screening).Error; err != nil {
			return fmt.Errorf("failed to create screening: %w", err)
		}

		var seat models.Seat
		if err := tx.
			Where("hall_id = ?", input.HallID).
			First(&seat).Error; err != nil {
			return fmt.Errorf("seat config not found: %w", err)
		}

		rows := seat.Row
		numbers := seat.Number

		seatScreenings := make(
			[]models.SeatScreening,
			0,
			rows*numbers,
		)

		for row := int64(1); row <= rows; row++ {
			for number := int64(1); number <= numbers; number++ {
				seatScreenings = append(seatScreenings, models.SeatScreening{
					HallID:      input.HallID,
					ScreeningID: screening.ID,
					Row:         row,
					Number:      number,
					Status:      "free",
				})
			}
		}

		if err := tx.Create(&seatScreenings).Error; err != nil {
			return fmt.Errorf("failed to create seat_screening: %w", err)
		}

		return nil
	})

	if err != nil {
		return models.Screening{}, err
	}

	return screening, nil
}

func (r *repository) Get(ctx context.Context, id int64) (models.Screening, error) {
	var screening models.Screening

	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&screening, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return screening, ErrNotFound
		}

		return screening, fmt.Errorf("failed to get Screening: %w", err)
	}

	return screening, nil
}

func (r *repository) Delete(ctx context.Context, id int64) (bool, error) {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Screening{})
	if res.Error != nil {
		return false, fmt.Errorf("failed to delete Screening: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return false, ErrNotFound
	}

	return true, nil
}

func (r *repository) ListByFilm(ctx context.Context, filmId int64) ([]models.Screening, error) {
	var screenings []models.Screening

	if err := r.db.WithContext(ctx).Where("film_id = ?", filmId).
		Find(&screenings).Error; err != nil {

		return nil, fmt.Errorf("failed to list screening by film: %w", err)
	}

	return screenings, nil
}

func (r *repository) ListByCinema(ctx context.Context, cinemaId int64) ([]models.Screening, error) {
	var screenings []models.Screening

	if err := r.db.WithContext(ctx).Where("cinema_id = ?", cinemaId).
		Find(&screenings).Error; err != nil {

		return nil, fmt.Errorf("failed to list screening by cinema: %w", err)
	}

	return screenings, nil
}

func (r *repository) ListByHall(ctx context.Context, hallId int64) ([]models.Screening, error) {
	var screenings []models.Screening

	if err := r.db.WithContext(ctx).Where("hall_id = ?", hallId).
		Find(&screenings).Error; err != nil {

		return nil, fmt.Errorf("failed to list screening by hall: %w", err)
	}

	return screenings, nil
}
