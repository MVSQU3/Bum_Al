package controller

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"xxx/models"
	"xxx/repo"
	"xxx/utils"

	"github.com/gin-gonic/gin"
)

type AlbumController struct {
	repo *repo.AlbumRepository
}

func NewAlbumController(repo *repo.AlbumRepository) *AlbumController {
	return &AlbumController{repo: repo}
}

func (ctrl *AlbumController) GetAllAlbums(c *gin.Context) {
	albums, err := ctrl.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, albums)
}

func (ctrl *AlbumController) GetAlbumsById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param(("id")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalide album id"})
		return
	}

	album, err := ctrl.repo.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, album)
}

func (ctrl *AlbumController) AddAlbums(c *gin.Context) {
	// var album models.Album
	// if err := c.BindJSON(&album); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// }

	file, err := c.FormFile("cover")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur lecture fichier"})
		return
	}

	var album models.Album
	album.Title = c.PostForm("title")
	album.Artist = c.PostForm("artist")
	album.Year, _ = strconv.Atoi(c.PostForm("year"))

	if file != nil {
		ctx := context.Background()
		filePath := "/tmp/" + file.Filename
		// Sauvegarder le fichier temporairement
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "erreur de sauvegarde fichier"})
			return
		}

		// Upload à Cloudinary
		coverURL, err := utils.UploadImage(ctx, filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur upload image"})
			return
		}
		album.Cover_url = coverURL
		defer os.Remove(filePath)

	}

	createdAlbum, err := ctrl.repo.Add(&album)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdAlbum)
}

func (ctrl *AlbumController) UpdateAlbums(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid album id"})
		return
	}

	// Récupérer le fichier depuis la requête multipart
	file, err := c.FormFile("cover")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur lecture fichier"})
		return
	}

	var album models.Album
	album.Title = c.PostForm("title")
	album.Artist = c.PostForm("artist")
	album.Year, _ = strconv.Atoi(c.PostForm("year"))

	// Si un fichier est fourni, l'uploader à Cloudinary
	if file != nil {
		ctx := context.Background()
		filePath := "/tmp/" + file.Filename
		// Sauvegarder le fichier temporairement
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "erreur de sauvegarde fichier"})
			return
		}

		// Upload à Cloudinary
		coverURL, err := utils.UploadImage(ctx, filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur upload image"})
			return
		}
		album.Cover_url = coverURL
	}

	updatedAlbum, err := ctrl.repo.Update(id, &album)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedAlbum)

}

func (ctrl *AlbumController) DeleteAlbums(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "album not found"})
		return
	}
	err = ctrl.repo.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "album supprimé"})
}

// L'appel fmt.Sprintf a des arguments mais pas de directives de formatage
// Le format fmt.Sprintf %s a un argument mis à jourAlbum d'un type incorrect *xxx/models.Album

/*func GetAlbumsHandler(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT id, title, artist, year FROM albums")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var a models.Album
		if err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Year); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("a =>", a)
		albums = append(albums, a)
	}

	c.JSON(http.StatusOK, albums)
}*/

/*func GetAlbumsById(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	var album models.Album
	err := db.QueryRow("SELECT id, title, artist, year FROM albums WHERE id = $1", id).Scan(&album.ID, &album.Title, &album.Artist, &album.Year)

	switch {
	case err == sql.ErrNoRows:
		c.JSON(http.StatusNotFound, gin.H{"Message": "Aucun albums trouver"})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	default:
		c.JSON(http.StatusOK, album)
	}
}*/

/*func AddAlbums(c *gin.Context, db *sql.DB) {
	var album models.Album

	if err := c.BindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	_, err := db.Exec("INSERT INTO albums (title, artist, year) VALUES ($1, $2, $3)", album.Title, album.Artist, album.Year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Album %s ajouté avec succès", album.Title)})
}*/

/*func DeleteAlbums(c *gin.Context, db *sql.DB) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM albums WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	affecedRow, _ := result.RowsAffected()
	if affecedRow == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Aucun album trouver"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Album supprimer avec succès"})
}*/

/*func UpdateAlbums(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var album models.Album

	if err := c.BindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("UPDATE albums SET title=$1, year=$2 WHERE id=$3", album.Title, album.Year, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	affecedRow, _ := result.RowsAffected()

	if affecedRow == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Échec de modification"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Album a été modifier avec succès"})
}*/
