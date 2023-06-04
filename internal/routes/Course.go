package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/handlers"
)

func CourseSubRoutes(course *echo.Group) {
	course.POST("", handlers.CreateCourse)
}
