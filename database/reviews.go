package database

import (
	"context"
	"database/sql"
	"time"
)

type ReviewPlaces interface {
	SaveReview(ctx context.Context, name, address, comment,
		username string, rating int) error

	GetReview(ctx context.Context, name string) ([]string, []string, []int, []string, error)

	DeleteReview(ctx context.Context, name string) error

	UpdateReview(ctx context.Context)
}

type ReviewPlacesImpl struct {
	DB *sql.DB
}

func NewReviewPlacesImp(db *sql.DB) *ReviewPlacesImpl {
	return &ReviewPlacesImpl{
		DB: db,
	}
}

func (r *ReviewPlacesImpl) SaveReview(ctx context.Context, name, address, comment,
	username string, rating int) error {
	now := time.Now()
	if _, err := r.DB.ExecContext(ctx,`INSERT INTO Reviews(name,address,userName,rating,comment,created_at) VALUES($1,$2,$3,$4,$5,$6)`,
		name, address, username, rating, comment, now); err != nil {
		return err
	}

	return nil
}

func (r *ReviewPlacesImpl) GetReview(ctx context.Context, name string) ([]string, []string, []int, []string, error) {
	rows, err := r.DB.QueryContext(ctx,`SELECT userName,rating,comment,created_at FROM Reviews WHERE name=$1`, name)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	defer rows.Close()

	var comments []string
	var usernames []string
	var ratings []int
	var times_created []string

	for rows.Next() {
		var comment string
		var username string
		var rating int64
		var created string

		if err = rows.Scan(&username, &rating, &comment, &created); err != nil {
			return nil, nil, nil, nil, err
		}

		comments = append(comments, comment)
		usernames = append(usernames, username)
		ratings = append(ratings, int(rating))
		times_created = append(times_created, created)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, nil, nil, err
	}

	return comments, usernames, ratings, times_created, nil
}
