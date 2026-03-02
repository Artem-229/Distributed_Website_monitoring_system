package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CheckJWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")

		if authorization == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "declined",
				"message": "user not authorized",
			})
			return
		}

		authorization = strings.TrimPrefix(authorization, "bearer ")
		j := jwt.MapClaims{}
		err := jwt.ParseWithClaims(authorization, &j, func(word *jwt.Token) (interface{}, error) {
			return []byte(word), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("UserID", j["userid"])
		c.Next()
	}
}
