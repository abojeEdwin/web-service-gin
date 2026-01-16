package repository

import (
	"database/sql"
	"example/web-service-gin/models"
)

type AlbumRepository struct {
	db *sql.DB
}

func NewAlbumRepository(db *sql.DB) *AlbumRepository {
	return &AlbumRepository{db: db}
}

func (r *AlbumRepository) GetAll() ([]models.Album, error) {
	var albums []models.Album

	rows, err := r.db.Query("SELECT * FROM album")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var alb models.Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price, &alb.Genre); err != nil {
			return nil, err
		}
		albums = append(albums, alb)
	}
	return albums, nil
}

func (r *AlbumRepository) Add(alb models.Album) (int64, error) {
	result, err := r.db.Exec("INSERT INTO album (title, artist, price, genre) VALUES (?, ?, ?, ?)", alb.Title, alb.Artist, alb.Price, alb.Genre)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *AlbumRepository) GetByID(id string) (models.Album, error) {
	var alb models.Album
	row := r.db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price, &alb.Genre); err != nil {
		return alb, err
	}
	return alb, nil
}
