package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/constants"
	"github.com/mjaliz/deviran/internal/initializers"
	"github.com/mjaliz/deviran/internal/message"
	"github.com/mjaliz/deviran/internal/models"
	"github.com/mjaliz/deviran/internal/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
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

func DeserializeUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var accessToken string
		authorization := c.Request().Header.Get("Authorization")
		cookieToken, err := c.Cookie("access_token")
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, message.StatusErrMessage(err.Error()))
		}

		if strings.HasPrefix(authorization, "Bearer ") {
			accessToken = strings.TrimPrefix(authorization, "Bearer ")
		} else if cookieToken.Value != "" {
			accessToken = cookieToken.Value
		}

		if accessToken == "" {
			return c.JSON(http.StatusUnauthorized, message.StatusErrMessage("You are not logged in"))
		}

		config, _ := initializers.LoadConfig(".")

		tokenClaims, err := utils.ValidateToken(accessToken, config.AccessTokenPublicKey)
		if err != nil {
			return c.JSON(http.StatusForbidden, message.StatusErrMessage(err.Error()))
		}

		ctx := context.TODO()
		userid, err := initializers.RedisClient.Get(ctx, tokenClaims.TokenUuid).Result()
		if err == redis.Nil {
			return c.JSON(http.StatusForbidden, message.StatusErrMessage("token is invalid or session has expired"))
		}

		var user models.User
		err = initializers.DB.First(&user, fmt.Sprintf("%s = ?", models.UserIdField), userid).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusForbidden, message.StatusErrMessage("the user belonging to this token no logger exists"))
		}

		c.Set(constants.EchoUserAttribute, models.FilterUserRecord(&user))
		c.Set(constants.EchoAccessTokenUuid, tokenClaims.TokenUuid)

		return next(c)
	}
}
