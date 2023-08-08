package database

import (
	"context"
	"database/sql"
)

type SavedPlaces interface {
	SavePlace(ctx context.Context, name, address string) error

	GetSavePlaces(ctx context.Context) ([]string, []string, error)

	DeletePlace(ctx context.Context, name string) error
}

type SavedPlacesImpl struct {
	DB *sql.DB
}

func NewSavedPlacesImpl(db *sql.DB) *SavedPlacesImpl {
	return &SavedPlacesImpl{
		DB: db,
	}
}

func (s *SavedPlacesImpl) SavePlace(ctx context.Context, name, address string) error {
	if _, err := s.DB.Exec(`INSERT INTO SavePlace(name,address) VALUES($1,$2)`, name, address); err != nil {
		return err
	}
	return nil
}

func (s *SavedPlacesImpl) GetSavePlaces(ctx context.Context) ([]string, []string, error) {
	rows, err := s.DB.Query(`SELECT name,address FROM SavePlace`)
	if err != nil {
		return nil, nil, err
	}

	var names, addresses []string
	for rows.Next() {
		var name, address string
		err = rows.Scan(&name, &address)
		if err != nil {
			return nil, nil, err
		}
		
		names = append(names, name)
		addresses = append(addresses, address)
	}

	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return names, addresses, nil
}

func (s *SavedPlacesImpl) DeletePlace(ctx context.Context, name string) error {
	return nil
}
