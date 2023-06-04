package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/handlers"
	"github.com/mjaliz/deviran/internal/middleware"
)

func UserSubRoutes(user *echo.Group) {
	user.POST("/sign_up", handlers.SignUp)
	user.POST("/sign_in", handlers.SignIn)
	user.GET("/logout", handlers.Logout, middleware.DeserializeUser)
}
