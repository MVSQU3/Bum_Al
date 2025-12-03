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

	c.JSON(http.StatusOK, gin.H{"update": updatedAlbum, "success": true})

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
