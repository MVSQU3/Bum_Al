package repo

import (
	"database/sql"
	"fmt"
	"xxx/models"
)

type AlbumRepository struct {
	db *sql.DB
}

func NewAlbumRepository(db *sql.DB) *AlbumRepository {
	return &AlbumRepository{db: db}
}

// GetAll - récupérer tous les albums
func (r *AlbumRepository) GetAll() ([]models.Album, error) {
	query := "SELECT id, title, artist, year, cover_url FROM albums"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		album := &models.Album{}
		err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Year, &album.Cover_url)
		if err != nil {
			return nil, err
		}
		albums = append(albums, *album)
	}
	return albums, nil
}

// GetById - récupérer un album
func (r *AlbumRepository) GetById(id int) (*models.Album, error) {
	query := "SELECT id, title, artist, year FROM albums WHERE id = $1"
	var album models.Album

	err := r.db.QueryRow(query, id).Scan(&album.ID, &album.Title, &album.Artist, &album.Year)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("album not found")
	}
	if err != nil {
		return nil, err
	}
	return &album, nil

}

// Add - ajouter un nouvel album
func (r *AlbumRepository) Add(album *models.Album) (*models.Album, error) {
	query := "INSERT INTO albums (title, artist, year, conver_url) VALUES ($1, $2, $3, $4) RETURNING id"
	err := r.db.QueryRow(query, album.Title, album.Artist, album.Year, album.Cover_url).Scan(&album.ID)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (r *AlbumRepository) Update(id int, album *models.Album) (*models.Album, error) {
	query := "UPDATE albums SET title = $1, artist = $2, year = $3, cover_url = $4  WHERE id=$5 RETURNING id, title, artist, year, cover_url"
	err := r.db.QueryRow(query, album.Title, album.Artist, album.Year, album.Cover_url, id).Scan(&album.ID, &album.Title, &album.Artist, &album.Year, &album.Cover_url)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("album not found")
	}
	if err != nil {
		return nil, err
	}
	return album, nil

}

// Delete - supprimer un album
func (r *AlbumRepository) Delete(id int) error {
	query := "DELETE FROM albums WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("album not found")
	}
	return nil
}
