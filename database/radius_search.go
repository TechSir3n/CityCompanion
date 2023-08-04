package database

import (
	"context"
	"database/sql"
	_"github.com/lib/pq"
)

type RadiusSearch interface {
	SaveRadiusSearch(ctx context.Context, radius float64) error

	GetRadiusSearch(ctx context.Context) (error, float64)

	UpdateRadiusSearch(ctx context.Context) (error)

	DeleteRadiusSearch(ctx context.Context) (error)
}

type RadiusSearchImpl struct {
	DB *sql.DB
}

func NewRadiusSearchImpl(db *sql.DB) *RadiusSearchImpl {
	return &RadiusSearchImpl{
		DB: db,
	}
}

func (r *RadiusSearchImpl) SaveRadiusSearch(ctx context.Context, radius float64) error {
	_, err := r.DB.Exec(`INSERT INTO SaveRadius(radius) VALUES($1)`, radius)
	if err != nil {
		return err
	}
	return nil
}


func (r *RadiusSearchImpl) GetRadiusSearch(ctx context.Context) (error, float64) {
	rows, err := r.DB.Query(`SELECT radius FROM SaveRadius`)
	if err != nil {
		return err, 0.0
	}

	defer rows.Close()

	var radius float64
	for rows.Next() {
		err := rows.Scan(&radius)
		if err != nil {
			return err, 0.0
		}
	}

	if err = rows.Err(); err != nil {
		return err, 0.0
	}

	return nil, radius
}


func (r *RadiusSearchImpl) UpdateRadiusSearch(ctx context.Context) (error) { 

	return nil
}