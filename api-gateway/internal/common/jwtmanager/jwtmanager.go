package jwtmanager

import (
	"errors"
	"fmt"
	"time"

	"github.com/ciscapello/api-gateway/internal/application/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type JwtManager struct {
	accsTokenExpires    time.Duration
	refreshTokenExpires time.Duration
	accsTokenSecret     string
	refreshTokenSecret  string

	logger *zap.Logger
}

func NewJwtManager(config *config.Config, logger *zap.Logger) *JwtManager {
	return &JwtManager{
		accsTokenExpires:    time.Duration(config.AccessTokenExpTime),
		refreshTokenExpires: time.Duration(config.RefreshTokenExpTime),
		accsTokenSecret:     config.AccessTokenSecret,
		refreshTokenSecret:  config.RefreshTokenSecret,

		logger: logger,
	}
}

func (j *JwtManager) Generate(uid uuid.UUID) (ReturnTokenType, error) {

	accessToken, err := j.genAccessToken(uid.String())
	if err != nil {
		j.logger.Error("failed to generate access token", zap.Error(err))
		return ReturnTokenType{}, err
	}

	refreshToken, err := j.genRefreshToken(uid.String())
	if err != nil {
		j.logger.Error("failed to generate refresh token", zap.Error(err))
		return ReturnTokenType{}, err
	}

	return ReturnTokenType{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *JwtManager) genAccessToken(uid string) (string, error) {
	fmt.Println(time.Now())
	fmt.Println(time.Now().Add(j.accsTokenExpires))
	fmt.Println(j.accsTokenExpires)
	fmt.Println(time.Now().Add(j.accsTokenExpires).Unix())

	claims := jwt.MapClaims{
		"id":  uid,
		"exp": time.Now().Add(j.accsTokenExpires).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(j.accsTokenSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (j *JwtManager) genRefreshToken(uid string) (string, error) {
	claims := jwt.MapClaims{
		"id":  uid,
		"exp": time.Now().Add(j.refreshTokenExpires).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := token.SignedString([]byte(j.refreshTokenSecret))
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (j *JwtManager) Verify(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.accsTokenSecret), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["id"].(string), nil
	}

	return "", errors.New("invalid token")
}

func (j *JwtManager) Parse() (string, error) {
	return "", nil
}

func (j *JwtManager) Refresh() (string, error) {
	return "", nil
}

func (j *JwtManager) Revoke() (string, error) {
	return "", nil
}

func (j *JwtManager) Validate() (string, error) {
	return "", nil
}
