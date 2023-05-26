package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/models"
	"github.com/mjaliz/deviran/internal/utils"
	"log"
	"net/http"
)

func (m *Repository) SignUp(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(user); err != nil {
		return err
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println("Hashing password failed")
	}
	user.Password = hashedPassword
	return c.JSON(http.StatusCreated, user)
}
