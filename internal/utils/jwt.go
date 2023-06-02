package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/constants"
	"os"
	"strings"
	"time"
)

type JwtPayload struct {
	jwt.RegisteredClaims
	ID     uuid.UUID `json:"id"`
	UserId int       `json:"user_id"`
}

var (
	errExtractingToken = errors.New("extracting token from header failed")
	errInvalidToken    = errors.New("token is invalid")
)

func newPayload(userId int) (JwtPayload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return JwtPayload{}, err
	}

	jwtPayload := JwtPayload{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.JWTValidDuration * time.Minute)),
		},
		ID:     tokenID,
		UserId: userId,
	}
	return jwtPayload, nil
}

func CreateToken(userID int) (string, error) {
	payload, err := newPayload(userID)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func verifyToken(token string) (*JwtPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errInvalidToken
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JwtPayload{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}
		return nil, errInvalidToken
	}

	payload, ok := jwtToken.Claims.(*JwtPayload)
	if !ok {
		return nil, errInvalidToken
	}

	return payload, nil
}

func extractToken(c echo.Context) string {
	bearToken := c.Request().Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func ExtractTokenMetadata(c echo.Context) (*JwtPayload, error) {
	token := extractToken(c)
	if token == "" {
		return nil, errExtractingToken
	}
	jwtPayload, err := verifyToken(token)
	if err != nil {
		return nil, err
	}
	return jwtPayload, nil
}
