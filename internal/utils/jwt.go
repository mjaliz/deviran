package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mjaliz/deviran/internal/constants"
	"github.com/mjaliz/deviran/internal/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// `iat`, `nbf`, `iss`, `sub` and `aud`.

type JwtPayload struct {
	ID        uuid.UUID `json:"id"`
	UserId    int       `json:"uid"`
	IssuedAt  time.Time `json:"iat"`
	ExpiredAt time.Time `json:"exp"`
}

type JWTMaker struct {
	secretKey string
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func NewPayload(userId int) (jwt.Claims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	claims := jwt.MapClaims{}
	claims["id"] = tokenID
	claims["uid"] = userId
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(constants.JWTValidDuration * time.Minute)

	return claims, nil
}

func NewJWTMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < constants.MinSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", constants.MinSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(userID int) (string, error) {
	payload, err := NewPayload(userID)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*JwtPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JwtPayload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func CreateToken(userId uint64) (*models.TokenDetails, error) {
	t := &models.TokenDetails{}
	t.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	t.AccessUuid = uuid.New().String()

	t.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	t.RefreshUuid = uuid.New().String()

	var err error
	//Creating Access TokenDetails
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = t.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = t.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	t.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh TokenDetails
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = t.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = t.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	t.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (m *Repository) CreateAuth(userId uint64, td *models.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := m.App.RedisClient.Set(m.App.RedisCtx, td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := m.App.RedisClient.Set(m.App.RedisCtx, td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
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

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (*models.AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &models.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func FetchAuth(authD *models.AccessDetails) (uint64, error) {
	userid, err := Repo.App.RedisClient.Get(Repo.App.RedisCtx, authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}
