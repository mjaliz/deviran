package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mjaliz/deviran/internal/models"
	"os"
	"strconv"
	"time"
)

func CreateToken(userId uint) (*models.TokenDetails, error) {
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

func (m *Repository) CreateAuth(userId uint, td *models.TokenDetails) error {
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
