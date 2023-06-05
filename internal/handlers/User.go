package handlers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/constants"
	"github.com/mjaliz/deviran/internal/initializers"
	"github.com/mjaliz/deviran/internal/input"
	"github.com/mjaliz/deviran/internal/message"
	"github.com/mjaliz/deviran/internal/models"
	"github.com/mjaliz/deviran/internal/output"
	"github.com/mjaliz/deviran/internal/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

func SignUp(c echo.Context) error {
	var payload input.SignUp

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, message.StatusBadRequestMessage(nil, err.Error()))
	}

	validateErrors := models.ValidateStruct(payload)
	if validateErrors != nil {
		return c.JSON(http.StatusBadRequest, message.StatusBadRequestMessage(validateErrors, ""))
	}

	if payload.Password != payload.PasswordConfirm {
		return c.JSON(http.StatusBadRequest, message.StatusBadRequestMessage(nil, "password do not match"))
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		log.Println("Hashing password failed", err.Error())
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}

	newUser := &models.User{
		Name:        payload.Name,
		Email:       strings.ToLower(payload.Email),
		Password:    hashedPassword,
		Photo:       &payload.Photo,
		PhoneNumber: payload.PhoneNumber,
	}

	err = initializers.DB.Create(&newUser).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.JSON(http.StatusConflict, message.StatusConflictMessage(err.Error()))
		}
		log.Println("creating user in db failed!", err.Error())
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}

	return c.JSON(http.StatusCreated, message.StatusOkMessage(models.FilterUserRecord(newUser), ""))
}

func SignIn(c echo.Context) error {
	var payload input.SignIn

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, message.StatusBadRequestMessage(nil, ""))
	}

	validateErrors := models.ValidateStruct(payload)
	if validateErrors != nil {
		return c.JSON(http.StatusBadRequest, message.StatusBadRequestMessage(validateErrors, ""))
	}

	var user models.User

	err := initializers.DB.First(&user, fmt.Sprintf("%s = ?", models.UserEmailField),
		strings.ToLower(payload.Email)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusForbidden, message.StatusForbiddenMessage("Invalid email or password"))
		}
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}

	err = utils.CompareHashAndPass(user.Password, payload.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}

	accessTokenDetails, err := utils.CreateToken(user.ID, initializers.Config.AccessTokenExpiresIn, initializers.Config.AccessTokenPrivateKey)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, message.StatusErrMessage(err.Error()))
	}

	refreshTokenDetails, err := utils.CreateToken(user.ID, initializers.Config.RefreshTokenExpiresIn, initializers.Config.RefreshTokenPrivateKey)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, message.StatusErrMessage(err.Error()))
	}

	now := time.Now()

	errAccess := initializers.RedisClient.Set(initializers.Ctx, accessTokenDetails.TokenUuid,
		user.ID, time.Unix(*accessTokenDetails.ExpiresIn, 0).Sub(now)).Err()
	if errAccess != nil {
		return c.JSON(http.StatusUnprocessableEntity, message.StatusErrMessage(errAccess.Error()))
	}

	errRefresh := initializers.RedisClient.Set(initializers.Ctx, refreshTokenDetails.TokenUuid,
		user.ID, time.Unix(*refreshTokenDetails.ExpiresIn, 0).Sub(now)).Err()
	if errRefresh != nil {
		return c.JSON(http.StatusUnprocessableEntity, message.StatusErrMessage(errRefresh.Error()))
	}

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    *accessTokenDetails.Token,
		Path:     "/",
		MaxAge:   initializers.Config.AccessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
		Domain:   "localhost",
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    *refreshTokenDetails.Token,
		Path:     "/",
		MaxAge:   initializers.Config.RefreshTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
		Domain:   "localhost",
	})

	c.SetCookie(&http.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   initializers.Config.AccessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: false,
		Domain:   "localhost",
	})

	return c.JSON(http.StatusOK, message.StatusOkMessage(output.SignIn{AccessToken: *accessTokenDetails.Token}, ""))
}

func RefreshAccessToken(c echo.Context) error {
	errMessage := "could not refresh access token"

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}

	if refreshToken.Value == "" {
		return c.JSON(http.StatusForbidden, message.StatusErrMessage(errMessage))
	}

	tokenClaims, err := utils.ValidateToken(refreshToken.Value, initializers.Config.RefreshTokenPublicKey)
	if err != nil {
		return c.JSON(http.StatusForbidden, message.StatusErrMessage(err.Error()))
	}

	userid, err := initializers.RedisClient.Get(initializers.Ctx, tokenClaims.TokenUuid).Result()
	if err == redis.Nil {
		return c.JSON(http.StatusForbidden, message.StatusErrMessage(errMessage))
	}

	var user models.User

	err = initializers.DB.First(&user, fmt.Sprintf("%s = ?", models.UserIdField), userid).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusForbidden, message.StatusErrMessage("the user belonging to this token no logger exists"))
		}
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}

	accessTokenDetails, err := utils.CreateToken(user.ID, initializers.Config.AccessTokenExpiresIn, initializers.Config.AccessTokenPrivateKey)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, message.StatusErrMessage(err.Error()))
	}

	now := time.Now()

	errAccess := initializers.RedisClient.Set(initializers.Ctx, accessTokenDetails.TokenUuid, user.ID,
		time.Unix(*accessTokenDetails.ExpiresIn, 0).Sub(now)).Err()
	if errAccess != nil {
		return c.JSON(http.StatusUnprocessableEntity, message.StatusErrMessage(errAccess.Error()))
	}

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    *accessTokenDetails.Token,
		Path:     "/",
		MaxAge:   initializers.Config.AccessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
		Domain:   "localhost",
	})

	c.SetCookie(&http.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   initializers.Config.AccessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: false,
		Domain:   "localhost",
	})

	return c.JSON(http.StatusOK, message.StatusOkMessage(accessTokenDetails.Token, ""))
}

func Logout(c echo.Context) error {
	errMessage := "Token is invalid or session has expired"

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		log.Println("reading refresh token from cookie failed", err.Error())
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}

	if refreshToken.Value == "" {
		return c.JSON(http.StatusForbidden, message.StatusErrMessage(errMessage))
	}

	tokenClaims, err := utils.ValidateToken(refreshToken.Value, initializers.Config.RefreshTokenPublicKey)
	if err != nil {
		return c.JSON(http.StatusForbidden, message.StatusErrMessage(err.Error()))
	}

	accessTokenUuid := c.Get(constants.EchoAccessTokenUuid).(string)
	_, err = initializers.RedisClient.Del(initializers.Ctx, tokenClaims.TokenUuid, accessTokenUuid).Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, message.StatusInternalServerErrorMessage())
	}

	expired := time.Now().Add(-time.Hour * 24)
	c.SetCookie(&http.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: expired,
	})
	c.SetCookie(&http.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Expires: expired,
	})
	c.SetCookie(&http.Cookie{
		Name:    "logged_in",
		Value:   "",
		Expires: expired,
	})
	return c.JSON(http.StatusOK, message.StatusOkMessage(nil, ""))
}

func GetMe(c echo.Context) error {
	user := c.Get(constants.EchoUserAttribute).(models.UserResponse)

	return c.JSON(http.StatusOK, message.StatusOkMessage(user, ""))
}
