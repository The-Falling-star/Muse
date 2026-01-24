package jwt

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ling/muse/config"
	"github.com/ling/muse/middleware"
)

// GenerateJWT 生成JWT Token
func GenerateJWT(userID int) (string, error) {
	cfg := config.Get()

	claims := &middleware.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token有效期24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Auth.JWTSecret))
}

// ParseJWT 解析JWT Token
func ParseJWT(tokenString string) (*middleware.JWTClaims, error) {
	cfg := config.Get()

	claims := &middleware.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Auth.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// GetUserId 获取用户Id
func GetUserId(ctx context.Context) int {
	if userId, ok := ctx.Value("userId").(int); ok && userId != 0 {
		return userId
	}
	return 0
}
