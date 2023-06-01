package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mjaliz/deviran/internal/constants"
	"net/http"
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
	ErrInvalidToken = errors.New("token is invalid")
)

func NewPayload(userId int) (JwtPayload, error) {
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
	payload, err := NewPayload(userID)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func VerifyToken(token string) (*JwtPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JwtPayload{}, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*JwtPayload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func ExtractTokenMetadata(r *http.Request) (*JwtPayload, error) {
	token := ExtractToken(r)
	jwtPayload, err := VerifyToken(token)
	if err != nil {
		return nil, err
	}
	return jwtPayload, nil
}
