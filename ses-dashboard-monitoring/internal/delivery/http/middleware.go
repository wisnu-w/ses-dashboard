package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		token, err := parser.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		if userID, ok := claims["user_id"]; ok {
			c.Set("user_id", userID)
		}
		if username, ok := claims["username"]; ok {
			c.Set("username", username)
		}

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if claims, exists := c.Get("claims"); exists {
			if jwtClaims, ok := claims.(jwt.MapClaims); ok {
				if role, exists := jwtClaims["role"]; exists && role == "admin" {
					c.Next()
					return
				}
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		c.Abort()
	}
}
