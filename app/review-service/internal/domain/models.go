package domain

import (
	"time"
)

type Review struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	FilmID    int       `json:"film_id" gorm:"not null;index"`
	UserID    int       `json:"user_id" gorm:"not null;index"`
	Title     string    `json:"title" gorm:"size:255"`
	Body      string    `json:"body" gorm:"type:text"`
	Rating    float64   `json:"rating" gorm:"type:decimal(3,1);check:rating >= 0 AND rating <= 5"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreateReviewRequest struct {
	FilmID int     `json:"film_id" binding:"required"`
	UserID int     `json:"user_id" binding:"required"`
	Title  string  `json:"title" binding:"required,max=255"`
	Body   string  `json:"body" binding:"required,min=10,max=5000"`
	Rating float64 `json:"rating" binding:"required,min=0,max=5"`
}

type UpdateReviewRequest struct {
	Title  string  `json:"title" binding:"omitempty,max=255"`
	Body   string  `json:"body" binding:"omitempty,min=10,max=5000"`
	Rating float64 `json:"rating" binding:"omitempty,min=0,max=5"`
}

type ListReviewsQuery struct {
	FilmID      int     `form:"film_id"`
	UserID      int     `form:"user_id"`
	MinRating   float64 `form:"min_rating"`
	MaxRating   float64 `form:"max_rating"`
	SortBy      string  `form:"sort_by"`
	SortOrder   string  `form:"sort_order"`
	Page        int     `form:"page" binding:"min=1"`
	PageSize    int     `form:"page_size" binding:"min=1,max=100"`
}

type FilmRatingStats struct {
	FilmID        int     `json:"film_id"`
	AverageRating float64 `json:"average_rating"`
	TotalReviews  int64   `json:"total_reviews"`
	Rating5Stars  int64   `json:"rating_5_stars"`
	Rating4Stars  int64   `json:"rating_4_stars"`
	Rating3Stars  int64   `json:"rating_3_stars"`
	Rating2Stars  int64   `json:"rating_2_stars"`
	Rating1Stars  int64   `json:"rating_1_stars"`
}

type ReviewWithUser struct {
	Review   Review   `json:"review"`
	UserName string   `json:"user_name"`
	UserID   int      `json:"user_id"`
}

func (Review) TableName() string {
	return "review"
}