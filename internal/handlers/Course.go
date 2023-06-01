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
	jwtPayload, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}

	course.UserId = jwtPayload.UserId
	return c.JSON(http.StatusOK, message.StatusOkMessage(course, ""))
}
