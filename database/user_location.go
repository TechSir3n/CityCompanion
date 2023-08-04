package database

import (
	"context"
	"database/sql"
)

type UserLocation interface {
	SaveUserLocation(ctx context.Context, latitude float64, longitude float64) error

	GetUserLocation(ctx context.Context) (error, float64, float64)

	UpdateUserLocation(ctx context.Context) error

	DeleteUserLocation(ctx context.Context) error
}

type UserLocationImpl struct {
	DB *sql.DB
}

func NewUserLocationImpl(db *sql.DB) *UserLocationImpl {
	return &UserLocationImpl{
		DB: db,
	}
}

func (u *UserLocationImpl) SaveUserLocation(ctx context.Context, latitude float64, longitude float64) error {
	_, err := u.DB.Exec(`INSERT INTO SaveLocation(latitude,longitude) VALUES($1,$2)`, latitude, longitude)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserLocationImpl) GetUserLocation(ctx context.Context) (error, float64, float64) {
	rows, err := u.DB.Query(`SELECT latitude,longitude FROM SaveLocation`)
	if err != nil {
		return err, 0.0, 0.0
	}

	defer rows.Close()

	var lati, longi float64
	for rows.Next() {
		err := rows.Scan(&lati, &longi)
		if err != nil {
			return err, 0.0, 0.0
		}
	}

	if err := rows.Err(); err != nil {
		return err, 0.0, 0.0
	}

	return nil, lati, longi
}

func (u *UserLocationImpl) UpdateUserLocation(ctx context.Context) error {
	return nil
}
