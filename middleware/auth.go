package middleware

import (
	"net/http"
	"soccer-team-api/model"
	"soccer-team-api/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header missing", nil)
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.ParseJWT(tokenStr)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token", nil)
			c.Abort()
			return
		}

		if !token.Valid {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token", nil)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims", nil)
			c.Abort()
			return
		}
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token has expired", nil)
			c.Abort()
			return
		}

		var user model.User
		if err := utils.DB.First(&user, claims["id"]).Error; err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User not found", nil)
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_email", user.Email)
		c.Next()
	}
}
