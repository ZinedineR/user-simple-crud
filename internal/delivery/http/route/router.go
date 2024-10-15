package route

import (
	"github.com/gin-gonic/gin"
	"user-simple-crud/internal/delivery/http"
	api "user-simple-crud/internal/delivery/http/middleware"
)

type Router struct {
	App            *gin.Engine
	UserHandler    *http.UserHTTPHandler
	AuthMiddleware *api.AuthMiddleware
}

func (h *Router) Setup() {
	h.App.Use(h.AuthMiddleware.ErrorHandler)
	guestApi := h.App.Group("/auth")
	{
		guestApi.POST("/register", h.UserHandler.Register)
		guestApi.POST("/login", h.UserHandler.Login)
	}
	coreApi := h.App.Group("")
	coreApi.Use(h.AuthMiddleware.JWTAuthentication)
	{
		userApi := coreApi.Group("/users")
		{
			userApi.POST("", h.UserHandler.Create)
			userApi.GET("", h.UserHandler.List)
			userApi.GET("/:id", h.UserHandler.FindOne)
			userApi.PUT("/:id", h.UserHandler.Update)
			userApi.DELETE("/:id", h.UserHandler.Delete)
		}
	}
}
