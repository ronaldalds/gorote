package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PayloadJwt struct {
	Token  string
	Claims JwtClaims
}
type JwtClaims struct {
	Sub         uint     `json:"sub"`
	Exp         int      `json:"exp"`
	Permissions []string `json:"permissions"`
	IsSuperUser bool     `json:"isSuperUser"`
	jwt.RegisteredClaims
}

type GenToken struct {
	Id          uint
	AppName     string
	Permissions []string
	IsSuperUser bool
	TimeZone    string
	JwtSecret   string
	Ttl         time.Duration
}

func GenerateToken(gen *GenToken) (string, error) {
	location, err := time.LoadLocation(gen.TimeZone)
	if err != nil {
		return "", fmt.Errorf("invalid timezone: %s", err.Error())
	}
	currentTime := time.Now().In(location)

	accessTokenExpirationTime := currentTime.Add(gen.Ttl)

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":         gen.Id,
		"iss":         gen.AppName,
		"permissions": gen.Permissions,
		"isSuperUser": gen.IsSuperUser,
		"iat":         currentTime.Unix(),
		"exp":         accessTokenExpirationTime.Unix(),
	})

	accessToken, err := accessClaims.SignedString([]byte(gen.JwtSecret))
	if err != nil {
		return "", fmt.Errorf("could not sign access token string %v", err.Error())
	}

	return accessToken, nil
}

func GetJwtHeaderPayload(auth, secret string) (*PayloadJwt, error) {
	// authHeader := ctx.Get("Authorization")
	tokenString := strings.Replace(auth, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JwtClaims{},
		func(t *jwt.Token) (any, error) {
			tokenSecret := secret
			return []byte(tokenSecret), nil
		},
	)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid jwt token")
	}

	tokenDone := token.Claims.(*JwtClaims)
	jwt := &PayloadJwt{
		Token:  tokenString,
		Claims: *tokenDone,
	}

	return jwt, nil
}
