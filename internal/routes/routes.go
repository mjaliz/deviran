package routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	customMiddleware "github.com/mjaliz/deviran/internal/middleware"
	"github.com/mjaliz/deviran/internal/models"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Routes() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &models.CustomValidator{Validator: validator.New()}
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	userGroup := e.Group("/user")
	UserSubRoutes(userGroup)

	courseGroup := e.Group("/course", customMiddleware.CheckAccessToken)
	CourseSubRoutes(courseGroup)
	e.Logger.Info(e.Start(":1323"))
}
