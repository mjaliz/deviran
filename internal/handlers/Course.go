package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/constants"
	"github.com/mjaliz/deviran/internal/message"
	"github.com/mjaliz/deviran/internal/models"
	"net/http"
)

func CreateCourse(c echo.Context) error {
	course := new(models.Course)
	if err := c.Bind(course); err != nil {
		return c.JSON(http.StatusBadRequest, message.StatusBadRequestMessage(nil, err.Error()))
	}
	err := models.ValidateStruct(course)
	if err != nil {
		return c.JSON(http.StatusBadRequest, message.StatusBadRequestMessage(err, ""))
	}

	user := c.Get(constants.EchoUserAttribute).(models.UserResponse)
	course.UserId = user.ID
	return c.JSON(http.StatusCreated, message.StatusOkMessage(course, ""))
}
