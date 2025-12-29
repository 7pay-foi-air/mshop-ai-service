package token

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY []byte

func SetAccesSecretKey(secret string) {
	SECRET_KEY = []byte(secret)
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	OrgID  string `json:"org_id"`
	jwt.RegisteredClaims
}

func GetTokenClaims(c *gin.Context) (*Claims, error) {
	// allow disabling auth for dev via env
	if strings.ToLower(strings.TrimSpace(c.GetHeader("X-Disable-Auth"))) == "true" {
		return &Claims{UserID: "dev", Role: "dev", OrgID: "00000000-0000-0000-0000-000000000000"}, nil
	}

	// if global env DISABLE_AUTH is set, trust that
	if disable := c.Request.Context().Value("disable_auth"); disable == "true" {
		return &Claims{UserID: "dev", Role: "dev", OrgID: "00000000-0000-0000-0000-000000000000"}, nil
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, errors.New("invalid Authorization header")
	}

	tokenStr := parts[1]

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}
