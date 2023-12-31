package database

import (
	"context"
	"database/sql"
)

type FavoritePlaces interface {
	SaveFavoritePlace(ctx context.Context, userID int64, name, address string) error

	GetFavoritePlaces(ctx context.Context, userID int64) ([]string, []string, error)

	DeleteFavoritePlaces(ctx context.Context, userID int64) error

	DeleteOneFavoritePlace(ctx context.Context, userID int64, name string) error
}

type FavoritePlacesImpl struct {
	DB *sql.DB
}

func NewFavoritePlacesImp(db *sql.DB) *FavoritePlacesImpl {
	return &FavoritePlacesImpl{
		DB: db,
	}
}

func (f *FavoritePlacesImpl) SaveFavoritePlace(ctx context.Context, userID int64, name, address string) error {
	if _, err := f.DB.ExecContext(ctx, `INSERT INTO  SavedFavoritePlace(userID,name,address) VALUES($1,$2,$3)`,
		userID, name, address); err != nil {
		return err
	}
	return nil
}

func (f *FavoritePlacesImpl) GetFavoritePlaces(ctx context.Context, userID int64) ([]string, []string, error) {
	rows, err := f.DB.QueryContext(ctx, `SELECT name,address FROM SavedFavoritePlace WHERE userID=$1`, userID)
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

func (f *FavoritePlacesImpl) DeleteFavoritePlaces(ctx context.Context, userID int64) error {
	if _, err := f.DB.ExecContext(ctx, `DELETE FROM SavedFavoritePlace WHERE userID=$1`,
		userID); err != nil {
		return err
	}
	return nil
}

func (f *FavoritePlacesImpl) DeleteOneFavoritePlace(ctx context.Context, userID int64, name string) error {
	if _, err := f.DB.ExecContext(ctx, `DELETE  FROM SavedFavoritePlace WHERE userID=$1 AND name=$2`,
		userID, name); err != nil {
		return err
	}
	return nil
}
