package main

import (
	"log"
	"time"

	"xxx/controller"
	"xxx/db"
	"xxx/middleware"
	"xxx/repo"

	"github.com/gin-contrib/cors"
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

	db := db.InitDb()
	defer db.Close()

	// Pour tester : 5 requêtes par minute avec burst de 2
	limiter := middleware.NewRateLimited(100/60.0, 50)

	// Initialiser repository
	albumRepo := repo.NewAlbumRepository(db)
	userRepo := repo.NewUserController(db)

	// Initialiser Controller
	albumCrtl := controller.NewAlbumController(albumRepo)
	userCttl := controller.NewUserController(userRepo)

	r := gin.Default()
	// 🔥 APPLIQUER le middleware CORS à toutes les routes
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Créer le groupe API avec JWT et rate limiting
	api := r.Group("/api")

	api.Use(middleware.ValidateJWT(), limiter.RateLimite())

	// ----- Routes avec rate limiting -----
	api.GET("/albums", albumCrtl.GetAllAlbums)
	api.GET("/albums/:id", albumCrtl.GetAlbumsById)
	api.POST("/albums", albumCrtl.AddAlbums)
	api.DELETE("/albums/:id", albumCrtl.DeleteAlbums)
	api.PUT("/albums/:id", albumCrtl.UpdateAlbums)
	api.POST("/auth/check", userCttl.CheckAuth)

	// ----- Routes publiques (sans rate limiting) -----
	r.POST("/api/register", userCttl.Register)
	r.POST("/api/login", userCttl.Login)
	r.POST("/api/logout", userCttl.Logout)

	r.Run(":8080")
}
