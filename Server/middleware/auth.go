package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// var SECRET = []byte("SECRET_KEY")

func ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token manquant"})
			c.Abort()
			return
		}

		tokenString, err := jwt.Parse(token, func(tokenString *jwt.Token) (any, error) {
			if _, ok := tokenString.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Méthode de signature invalide")
			}
			return os.Getenv("JWT_SECRET"), nil
		})
		if err != nil || !tokenString.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}

		if claims, ok := tokenString.Claims.(jwt.MapClaims); ok && tokenString.Valid {
			c.Set("userEmail", claims["email"])
			c.Set("userExp", claims["exp"])
		}

		c.Next()
	}
}
