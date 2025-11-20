package controller

import (
	"net/http"

	"xxx/middleware"
	"xxx/models"
	"xxx/repo"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	repo *repo.UserRepository
}

func NewUserController(repo *repo.UserRepository) *UserController {
	return &UserController{repo: repo}
}

func (ctrl *UserController) Login(c *gin.Context) {
	var input models.Input

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.repo.SignIn(&input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mot de passe incorrect"})
		return
	}

	token, err := middleware.GenerateJWT(input.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Connexion réussie", "token": token})
}

func (ctrl *UserController) Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registed, err := ctrl.repo.SignUp(&user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	registed.Password = string(hash)

	token, err := middleware.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("token", token, 3600, "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{"message": "Inscription réussie", "token": token})
}

func (ctrl *UserController) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Déconnexion réussite"})
}

/*func RegisterHandler(c *gin.Context, db *sql.DB) {
	var u models.User

	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var count int
	err := db.QueryRow("SELECT 1 FROM users WHERE email = $1", u.Email).Scan(&count)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Cet utilisateur existe déjà"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	_, err = db.Exec("INSERT INTO users (fullname, email, password) VALUES ($1, $2, $3)", u.FullName, u.Email, hash)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := middleware.GenerateJWT(u.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Inscription réussie", "token": token})
}*/

/*func LoginHandler(c *gin.Context, db *sql.DB) {
	var input models.Input

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbEmail, dbPassword string
	row := db.QueryRow("SELECT email, password FROM users WHERE email = $1", input.Email)
	err := row.Scan(&dbEmail, &dbPassword)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur introuvable"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mot de passe incorrect"})
		return
	}

	token, err := middleware.GenerateJWT(input.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("token", token, 3600, "/", "", false, true)


	c.JSON(http.StatusOK, gin.H{"message": "Connexion réussie", "token": token})

}*/
