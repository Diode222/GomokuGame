package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte("webapp"), nil
			})
		if err == nil && token != nil {
			if token.Valid {
				c.Next()
			} else {
				c.String(http.StatusUnauthorized, "Token is not valid")
			}
		} else {
			c.String(http.StatusUnauthorized, "Unauthorized access to this resource")
		}
	}
}
