package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/message"
	"github.com/mjaliz/deviran/internal/models"
	"github.com/mjaliz/deviran/internal/utils"
	"net/http"
)

func (m *Repository) CreateCourse(c echo.Context) error {
	course := new(models.Course)
	if err := c.Bind(course); err != nil {
		return c.JSON(http.StatusBadRequest, message.StatusBadRequestMessage(nil, err.Error()))
	}
	if err := c.Validate(course); err != nil {
		return err
	}
	tokenAuth, err := utils.ExtractTokenMetadata(c.Request())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	userId, err := utils.FetchAuth(tokenAuth)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	course.UserId = userId
	return c.JSON(http.StatusOK, message.StatusOkMessage(course, ""))
}
