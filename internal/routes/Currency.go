package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/handlers"
)

func CurrencySubRoutes(course *echo.Group) {
	course.GET("", handlers.GetCurrencies)
	course.GET("/:currency_id", handlers.GetCurrency)
}
