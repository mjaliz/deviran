package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/message"
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
	var userDB models.User
	result := m.App.DB.Debug().Where(&models.User{Email: user.Email}).First(&userDB)
	if result.Error != nil {
		log.Println("finding user in database failed", result.Error.Error())
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}
	if userDB.Email != "" {
		return c.JSON(http.StatusBadRequest,
			message.StatusBadRequestMessage(nil, "this email is already exits!"))
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println("Hashing password failed", err.Error())
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}
	user.Password = hashedPassword
	result = m.App.DB.Create(&user)
	if result.Error != nil {
		log.Println("creating user in db failed!", result.Error.Error())
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}
	return c.JSON(http.StatusCreated, message.StatusOkMessage(user, ""))
}
