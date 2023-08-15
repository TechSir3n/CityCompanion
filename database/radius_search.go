package database

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

type RadiusSearch interface {
	SaveRadiusSearch(ctx context.Context, userID int64, radius float64) error

	GetRadiusSearch(ctx context.Context, userID int64) (error, float64)

	UpdateRadiusSearch(ctx context.Context, userID int64, radius float64) error

	DeleteRadiusSearch(ctx context.Context, userID int64, radius float64) error
}

type RadiusSearchImpl struct {
	DB *sql.DB
}

func NewRadiusSearchImpl(db *sql.DB) *RadiusSearchImpl {
	return &RadiusSearchImpl{
		DB: db,
	}
}

func (r *RadiusSearchImpl) SaveRadiusSearch(ctx context.Context, userID int64, radius float64) error {
	if _, err := r.DB.ExecContext(ctx,`INSERT INTO Radius(userID,radius) VALUES($1,$2)`, userID, radius); err != nil {
		return err
	}
	return nil
}

func (r *RadiusSearchImpl) GetRadiusSearch(ctx context.Context, userID int64) (error, float64) {
	var radius float64
	err := r.DB.QueryRowContext(ctx,`SELECT radius FROM Radius WHERE userID=$1`, userID).Scan(&radius)
	if err != nil {
		return err, 0.0
	}

	return nil, radius
}

func (r *RadiusSearchImpl) UpdateRadiusSearch(ctx context.Context, userID int64, radius float64) error {
	fmt.Println("update function here")
	if _, err := r.DB.ExecContext(ctx,`UPDATE Radius SET radius=$1 WHERE userID=$2`, radius, userID); err != nil {
		return err
	}
	return nil
}

func (r *RadiusSearchImpl) DeleteRadiusSearch(ctx context.Context, userID int64, radius float64) error {
	if _, err := r.DB.ExecContext(ctx,`DELETE FROM Radius WHERE userID=$1 AND radius=$2`, userID, radius); err != nil {
		return err
	}
	return nil
}
