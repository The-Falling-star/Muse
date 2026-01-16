package middleware

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/config"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthInterceptor JWT认证中间件
func AuthInterceptor(skipProcedures []string) connect.UnaryInterceptorFunc {
	skipMap := make(map[string]bool)
	for _, proc := range skipProcedures {
		skipMap[proc] = true
	}

	cfg := config.Get()

	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// 跳过无需认证的接口
			if skipMap[req.Spec().Procedure] {
				return next(ctx, req)
			}

			// 从 Header 中获取 Authorization
			auth := req.Header().Get("Authorization")
			if auth == "" {
				return nil, errs.NewStandard(connect.CodeUnauthenticated, "请先登录")
			}

			// 解析 Bearer Token
			auth = strings.TrimPrefix(auth, "Bearer ")
			// 验证 JWT Token
			claims := &JWTClaims{}
			token, err := jwt.ParseWithClaims(auth, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.Auth.JWTSecret), nil
			})

			if err != nil {
				return nil, errs.NewStandard(connect.CodeUnauthenticated, "解析Token失败")
			}

			if !token.Valid {
				return nil, errs.NewStandard(connect.CodeUnauthenticated, "Token无效或已过期")
			}

			// 将用户ID存入 context
			ctx = context.WithValue(ctx, "userId", claims.UserID)

			return next(ctx, req)
		}
	}
}
