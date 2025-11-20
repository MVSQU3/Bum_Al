package main

import (
	"log"
	"os"

	// "database/sql"
	"xxx/controller"
	"xxx/db"
	"xxx/repo"

	"xxx/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Charger les variables d'environnement depuis .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Aucun fichier .env trouvé, utilisation des variables d'environnement système")
	}

	// Vérifier que CLOUDINARY_URL est défini
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	if cloudinaryURL == "" {
		log.Fatal("CLOUDINARY_URL n'est pas définie")
	}
	db := db.InitDb()
	defer db.Close()

	// Initialiser repository
	albumRepo := repo.NewAlbumRepository(db)
	userRepo := repo.NewUserController(db)

	// Initialiser Controller
	albumCrtl := controller.NewAlbumController(albumRepo)
	userCttl := controller.NewUserController(userRepo)

	r := gin.Default()

	api := r.Group("/api", middleware.ValidateJWT())
	// api.Use(middleware.ValidateJWT())

	// ----- Albums Routes -----
	r.GET("/api/albums", albumCrtl.GetAllAlbums)

	r.GET("/api/albums/:id", albumCrtl.GetAlbumsById)

	api.POST("/albums/", albumCrtl.AddAlbums)

	api.DELETE("/albums/:id", albumCrtl.DeleteAlbums)

	api.PUT("/albums/:id", albumCrtl.UpdateAlbums)

	// ---- User routes -----
	// r.POST("/api/register", func(c *gin.Context) {
	// 	controller.RegisterHandler(c, db)
	// })

	r.POST("/api/register", userCttl.Register)
	r.POST("/api/login", userCttl.Login)
	r.POST("/api/logout", userCttl.Logout)

	r.Run(":8080")
}
