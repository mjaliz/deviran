package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mjaliz/deviran/internal/models"
	"os"
	"time"
)

func CreateToken(userid uint64) (*models.Token, error) {
	t := &models.Token{}
	t.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	t.AccessUuid = uuid.New().String()

	t.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	t.RefreshUuid = uuid.New().String()

	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = t.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = t.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	t.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = t.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = t.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	t.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return t, nil
}
