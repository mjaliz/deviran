package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	_ "github.com/mjaliz/deviran/docs"
	"github.com/mjaliz/deviran/models"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

// @title Deviran API
// @version 1.0
// @description This is the server of Deviran platform
// @license.name Apache 2.0

func main() {
	e := echo.New()
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
