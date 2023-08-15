package database

import (
	"context"
	"database/sql"
)

type UserLocation interface {
	SaveUserLocation(ctx context.Context, userID int64, latitude, longitude float64) error

	GetUserLocation(ctx context.Context, userID int64) (error, float64, float64)

	UpdateUserLocation(ctx context.Context, userID int64, latitude, longitude float64) error

	DeleteUserLocation(ctx context.Context, userID int64) error
}

type UserLocationImpl struct {
	DB *sql.DB
}

func NewUserLocationImpl(db *sql.DB) *UserLocationImpl {
	return &UserLocationImpl{
		DB: db,
	}
}

func (u *UserLocationImpl) SaveUserLocation(ctx context.Context, userID int64, latitude, longitude float64) error {
	if _, err := u.DB.ExecContext(ctx,`INSERT INTO UserLocation(userID,latitude,longitude) VALUES($1,$2,$3)`,
		userID, latitude, longitude); err != nil {
		return err
	}
	return nil
}

func (u *UserLocationImpl) GetUserLocation(ctx context.Context, userID int64) (error, float64, float64) {
	var latitude, longitude float64
	err := u.DB.QueryRowContext(ctx,`SELECT latitude,longitude FROM UserLocation WHERE userID =$1`,
		userID).Scan(&latitude, &longitude)
	if err != nil {
		return err, 0.0, 0.0
	}
	return nil, latitude, longitude
}

func (u *UserLocationImpl) UpdateUserLocation(ctx context.Context, userID int64, latitude, longitude float64) error {
	if _, err := u.DB.ExecContext(ctx,`UPDATE UserLocation SET latitude=$1,longitude=$2 WHERE userID=$3`, latitude,
		longitude, userID); err != nil {
		return err
	}
	return nil
}

func (u *UserLocationImpl) DeleteUserLocation(ctx context.Context, userID int64) error {
	if _, err := u.DB.ExecContext(ctx,`DELETE FROM UserLocation WHERE userID=$1`, userID); err != nil {
		return err
	}
	return nil
}
