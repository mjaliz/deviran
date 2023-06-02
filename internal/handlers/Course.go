package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/constants"
	"github.com/mjaliz/deviran/internal/message"
	"github.com/mjaliz/deviran/internal/models"
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

	course.UserId = c.Get(constants.EchoUserIDAttribute).(int)
	return c.JSON(http.StatusOK, message.StatusOkMessage(course, ""))
}
