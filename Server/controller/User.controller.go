package controller

import (
	// "log"
	"fmt"
	"net/http"

	"xxx/models"
	"xxx/repo"
	"xxx/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	repo *repo.UserRepository
}

func NewUserController(repo *repo.UserRepository) *UserController {
	return &UserController{repo: repo}
}

func (ctrl *UserController) Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash AVANT d'envoyer en DB
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)

	// Maintenant seulement on sauvegarde
	_, err := ctrl.repo.SignUp(&user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 3600, "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{"message": "Inscription réussie", "user": user.Email})
}

func (ctrl *UserController) Login(c *gin.Context) {
	var input models.Input

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var password = input.Password
	// log.Printf("log de input: %+v", password)
	user, err := ctrl.repo.SignIn(&input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	// log.Printf("user connecté: %+v", user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mot de passe incorrect"})
		return
	}

	token, err := utils.GenerateJWT(input.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Connexion réussie", "user": user.Email})
}

func (ctrl *UserController) Logout(c *gin.Context) {
	// Afficher les cookies avant suppression
	if cookie, err := c.Cookie("token"); err == nil {
		fmt.Printf("Cookie token avant suppression: %s\n", cookie)
	}

	// Supprimer le cookie
	c.SetCookie("token", "", -1, "/", "", false, true)

	// Vérifier après suppression
	if cookie, err := c.Cookie("token"); err != nil {
		fmt.Println("Cookie token supprimé avec succès")
	} else {
		fmt.Printf("Cookie token toujours présent: %s\n", cookie)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Déconnexion réussie"})
}

func (ctrl *UserController) CheckAuth(c *gin.Context) {
	// Le middleware ValidateJWT a déjà été exécuté et a set les infos dans le contexte
	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"authenticated": false,
			"message":       "Non authentifié",
		})
		return
	}

	// userExp, _ := c.Get("userExp")

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"user": gin.H{
			"email": userEmail,
			// "exp":   userExp,
		},
		"message": "Authentifié avec succès",
	})
}