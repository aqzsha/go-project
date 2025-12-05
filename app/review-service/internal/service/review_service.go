package service

import (
	"context"
	"errors"
	"fmt"

	"review-service/internal/domain"
	"review-service/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrReviewNotFound      = errors.New("review not found")
	ErrDuplicateReview     = errors.New("user has already reviewed this film")
	ErrUnauthorized        = errors.New("unauthorized to modify this review")
	ErrInvalidRating       = errors.New("rating must be between 0 and 5")
	ErrFilmNotFound        = errors.New("film not found")
	ErrUserNotVerified     = errors.New("user must watch the film before reviewing")
)

type ReviewService interface {
	CreateReview(ctx context.Context, req domain.CreateReviewRequest) (*domain.Review, error)
	GetReviewByID(ctx context.Context, id int) (*domain.Review, error)
	GetFilmReviews(ctx context.Context, filmID int, query domain.ListReviewsQuery) ([]domain.Review, int64, error)
	GetUserReviews(ctx context.Context, userID int, query domain.ListReviewsQuery) ([]domain.Review, int64, error)
	UpdateReview(ctx context.Context, reviewID int, userID int, req domain.UpdateReviewRequest) (*domain.Review, error)
	DeleteReview(ctx context.Context, reviewID int, userID int) error
	GetFilmRatingStats(ctx context.Context, filmID int) (*domain.FilmRatingStats, error)
}

type reviewService struct {
	reviewRepo      repository.ReviewRepository
	filmClient      FilmClient
	ticketClient    TicketClient
}

func NewReviewService(
	reviewRepo repository.ReviewRepository,
	filmClient FilmClient,
	ticketClient TicketClient,
) ReviewService {
	return &reviewService{
		reviewRepo:   reviewRepo,
		filmClient:   filmClient,
		ticketClient: ticketClient,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, req domain.CreateReviewRequest) (*domain.Review, error) {
	if req.Rating < 0 || req.Rating > 5 {
		return nil, ErrInvalidRating
	}

	exists, err := s.filmClient.FilmExists(ctx, req.FilmID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify film: %w", err)
	}
	if !exists {
		return nil, ErrFilmNotFound
	}

	existingReview, err := s.reviewRepo.GetUserReviewForFilm(ctx, req.UserID, req.FilmID)
	if err == nil && existingReview != nil {
		return nil, ErrDuplicateReview
	}

	hasWatched, err := s.ticketClient.UserHasWatchedFilm(ctx, req.UserID, req.FilmID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ticket: %w", err)
	}
	if !hasWatched {
		return nil, ErrUserNotVerified
	}

	review := &domain.Review{
		FilmID: req.FilmID,
		UserID: req.UserID,
		Title:  req.Title,
		Body:   req.Body,
		Rating: req.Rating,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	go s.updateFilmRating(context.Background(), req.FilmID)

	return review, nil
}

func (s *reviewService) GetReviewByID(ctx context.Context, id int) (*domain.Review, error) {
	review, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReviewNotFound
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}
	return review, nil
}

func (s *reviewService) GetFilmReviews(ctx context.Context, filmID int, query domain.ListReviewsQuery) ([]domain.Review, int64, error) {
	query.FilmID = filmID
	reviews, total, err := s.reviewRepo.GetByFilmID(ctx, filmID, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get film reviews: %w", err)
	}
	return reviews, total, nil
}

func (s *reviewService) GetUserReviews(ctx context.Context, userID int, query domain.ListReviewsQuery) ([]domain.Review, int64, error) {
	query.UserID = userID
	reviews, total, err := s.reviewRepo.GetByUserID(ctx, userID, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user reviews: %w", err)
	}
	return reviews, total, nil
}

func (s *reviewService) UpdateReview(ctx context.Context, reviewID int, userID int, req domain.UpdateReviewRequest) (*domain.Review, error) {
	review, err := s.GetReviewByID(ctx, reviewID)
	if err != nil {
		return nil, err
	}

	if review.UserID != userID {
		return nil, ErrUnauthorized
	}

	if req.Title != "" {
		review.Title = req.Title
	}
	if req.Body != "" {
		review.Body = req.Body
	}
	if req.Rating > 0 {
		if req.Rating < 0 || req.Rating > 5 {
			return nil, ErrInvalidRating
		}
		review.Rating = req.Rating
	}

	if err := s.reviewRepo.Update(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	go s.updateFilmRating(context.Background(), review.FilmID)

	return review, nil
}

func (s *reviewService) DeleteReview(ctx context.Context, reviewID int, userID int) error {
	review, err := s.GetReviewByID(ctx, reviewID)
	if err != nil {
		return err
	}

	if review.UserID != userID {
		return ErrUnauthorized
	}

	if err := s.reviewRepo.Delete(ctx, reviewID); err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	go s.updateFilmRating(context.Background(), review.FilmID)

	return nil
}

func (s *reviewService) GetFilmRatingStats(ctx context.Context, filmID int) (*domain.FilmRatingStats, error) {
	stats, err := s.reviewRepo.GetFilmRatingStats(ctx, filmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rating stats: %w", err)
	}
	return stats, nil
}

func (s *reviewService) updateFilmRating(ctx context.Context, filmID int) {
	avgRating, err := s.reviewRepo.GetAverageRating(ctx, filmID)
	if err != nil {
		return
	}

	s.filmClient.UpdateFilmRating(ctx, filmID, avgRating)
}

type FilmClient interface {
	FilmExists(ctx context.Context, filmID int) (bool, error)
	UpdateFilmRating(ctx context.Context, filmID int, rating float64) error
}

type TicketClient interface {
	UserHasWatchedFilm(ctx context.Context, userID int, filmID int) (bool, error)
}