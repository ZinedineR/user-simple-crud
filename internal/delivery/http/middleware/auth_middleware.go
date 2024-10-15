package api

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"strings"
	"user-simple-crud/pkg/signature"
)

type AuthMiddleware struct {
	Middleware
	signaturer signature.Signaturer
}

func NewAuthMiddleware(signaturer signature.Signaturer) *AuthMiddleware {
	return &AuthMiddleware{signaturer: signaturer}
}

func (m *AuthMiddleware) JWTAuthentication(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	authFields := strings.Fields(authHeader)
	if len(authFields) != 2 || strings.ToLower(authFields[0]) != "bearer" {
		m.UnauthorizedJSON(c, "Invalid token")
		return
	}
	token := authFields[1]

	res, exception := m.signaturer.JWTCheck(token)
	if exception != nil {
		m.ExceptionJSON(c, exception)
		return
	}

	c.Set("username", res.Username)
	c.Set("access_token", res.Token)

	c.Next()
}

func (m *AuthMiddleware) ErrorHandler(c *gin.Context) {

	defer func() {
		if err0 := recover(); err0 != nil {
			slog.Any("error", err0)
			m.InternalErrorJSON(c, "Request is halted unexpectedly, please contact the administrator.", err0)
		}
	}()
	c.Next()
}
