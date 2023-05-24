package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/handlers"
)

func UserSubRoutes(user *echo.Group) {
	user.POST("/sign_up", handlers.Repo.SignUp)
}
