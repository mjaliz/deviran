package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/constants"
	"github.com/mjaliz/deviran/internal/message"
	"github.com/mjaliz/deviran/internal/utils"
	"log"
	"net/http"
)

func CheckAccessToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwtPayload, err := utils.ExtractTokenMetadata(c)
		if err != nil {
			log.Println("CheckAccessToken middleware failed", err.Error())
			return c.JSON(http.StatusUnauthorized, message.StatusUnauthorizedMessage(""))
		}
		c.Set(constants.EchoUserIDAttribute, jwtPayload.UserId)
		return next(c)
	}
}
