package handlers

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/message"
	"github.com/mjaliz/deviran/internal/models"
	"github.com/mjaliz/deviran/internal/utils"
	"gorm.io/gorm"
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
	err := m.App.DB.Where(&models.User{Email: user.Email}).First(&userDB).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("finding user in database failed", err.Error())
			return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
		}
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
	err = m.App.DB.Create(&user).Error
	if err != nil {
		log.Println("creating user in db failed!", err.Error())
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}

	ts, err := utils.CreateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	user.AccessToken = ts
	return c.JSON(http.StatusCreated, message.StatusOkMessage(user, ""))
}
