package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/message"
	"github.com/mjaliz/deviran/internal/utils"
	"net/http"
)

func CheckAccessToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := utils.ExtractToken(c)
		if token == "" {
			return c.JSON(http.StatusUnauthorized, message.StatusUnauthorizedMessage("not authorized"))
		}
		c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
		return next(c)
	}
}
