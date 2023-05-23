package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mjaliz/deviran/internal/config"
	"github.com/mjaliz/deviran/internal/models"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func routes(app *config.AppConfig) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Validator = &models.CustomValidator{Validator: validator.New()}
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/sign_up", func(c echo.Context) error {
		user := new(models.User)
		if err := c.Bind(user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := c.Validate(user); err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, user)
	})
	e.Logger.Info(e.Start(":1323"))
}
