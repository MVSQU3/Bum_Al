package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// Tes imports locaux (je suppose qu'ils sont corrects)
	"xxx/controller" // ⚠️ Vérifie que le chemin du module est bon
	"xxx/db"
	"xxx/middleware"
	"xxx/repo"
)

// --- 1. EMBED DU FRONTEND ---
// On charge tout le dossier "dist" dans la mémoire du programme
//
//go:embed dist/*
var staticFiles embed.FS

func main() {
	// Charger les variables d'environnement (localement)
	// Sur Skybot, ça échouera silencieusement et prendra les Variables du Panel, c'est parfait.
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: Pas de fichier .env, utilisation des variables système (Prod)")
	}

	dbConn := db.InitDb()
	defer dbConn.Close()

	// Initialiser repository & controller
	albumRepo := repo.NewAlbumRepository(dbConn)
	userRepo := repo.NewUserController(dbConn)
	albumCrtl := controller.NewAlbumController(albumRepo)
	userCttl := controller.NewUserController(userRepo)

	limiter := middleware.NewRateLimited(100/60.0, 50)

	// Mode Release pour la prod
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// --- 2. CONFIGURATION FRONTEND (Fichiers Statiques) ---
	// On récupère le sous-dossier "dist"
	distFS, _ := fs.Sub(staticFiles, "dist")

	// On sert les fichiers statiques (JS, CSS, Images qui sont dans dist/assets)
	r.StaticFS("/assets", http.FS(distFS))

	// Middleware CORS
	// Note : En mode "Fichier Unique", le frontend et le backend sont sur le même domaine.
	// CORS est moins critique, mais on le garde pour la sécurité ou les accès externes.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://ton-site-skybot.tech"}, // Ajoute ton URL de prod si besoin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// --- 3. ROUTING API ---
	api := r.Group("/api")
	api.Use(middleware.ValidateJWT(), limiter.RateLimite())

	// Routes protégées
	api.GET("/albums", albumCrtl.GetAllAlbums)
	api.GET("/albums/:id", albumCrtl.GetAlbumsById)
	api.POST("/albums", albumCrtl.AddAlbums)
	api.DELETE("/albums/:id", albumCrtl.DeleteAlbums)
	api.PUT("/albums/:id", albumCrtl.UpdateAlbums)
	api.POST("/auth/check", userCttl.CheckAuth)

	// Routes publiques API
	r.POST("/api/register", userCttl.Register)
	r.POST("/api/login", userCttl.Login)
	r.POST("/api/logout", userCttl.Logout)

	// --- 4. ROUTING SPA (React) ---
	// C'est ici qu'on dit : "Si ce n'est pas une route API, renvoie le HTML de React"
	r.NoRoute(func(c *gin.Context) {
		// Si l'utilisateur demande une page (ex: /login, /dashboard) -> on envoie index.html
		// Si c'est une api qui n'existe pas (ex: /api/truc-bizarre) -> on laisse le 404 JSON par défaut ou on gère
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.FileFromFS("index.html", http.FS(distFS))
		} else {
			c.JSON(404, gin.H{"error": "Route API non trouvée"})
		}
	})

	// --- 5. DÉMARRAGE SKYBOT ---
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Fallback local
	}

	log.Println("🚀 Serveur démarré sur le port :" + port)
	r.Run(":" + port)
}
