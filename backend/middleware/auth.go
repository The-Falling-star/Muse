package middleware

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ling/muse/common/constant"
	"github.com/ling/muse/common/errs"
	"github.com/ling/muse/config"
	log "github.com/sirupsen/logrus"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthInterceptor 认证拦截器
type AuthInterceptor struct {
	skipMap map[string]struct{}
}

// NewAuthInterceptor 创建认证拦截器实例
func NewAuthInterceptor(skipProcedures []string) *AuthInterceptor {
	skipMap := make(map[string]struct{})
	for _, proc := range skipProcedures {
		skipMap[proc] = struct{}{}
	}
	return &AuthInterceptor{skipMap: skipMap}
}

// WrapUnary 包装一元调用的认证拦截器
func (a *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		// 跳过无需认证的接口
		if config.Get().Auth.SkipAuth {
			ctx = context.WithValue(ctx, constant.UserIDKey, config.Get().Auth.AdminUserId)
			return next(ctx, req)
		}
		if _, exist := a.skipMap[req.Spec().Procedure]; exist {
			return next(ctx, req)
		}
		if ctxWithUserId, err := Auth(ctx, req.Header().Get("Authorization")); err != nil {
			log.Errorf("auth failed: %v", err)
			return nil, err
		} else {
			return next(ctxWithUserId, req)
		}
	}
}

// WrapStreamingClient 包装流式客户端调用的认证拦截器
func (a *AuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

// WrapStreamingHandler 包装流式处理调用的认证拦截器
func (a *AuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, req connect.StreamingHandlerConn) error {
		// 跳过无需认证的接口
		if config.Get().Auth.SkipAuth {
			ctx = context.WithValue(ctx, constant.UserIDKey, config.Get().Auth.AdminUserId)
			return next(ctx, req)
		}
		if _, exist := a.skipMap[req.Spec().Procedure]; exist {
			return next(ctx, req)
		}
		if ctxWithUserId, err := Auth(ctx, req.RequestHeader().Get("Authorization")); err != nil {
			log.Errorf("auth failed: %v", err)
			return err
		} else {
			return next(ctxWithUserId, req)
		}
	}
}

// Auth 验证JWT令牌并提取用户ID
func Auth(ctx context.Context, accessToken string) (context.Context, error) {
	// 解析 Bearer Token
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	// 验证 JWT Token
	claims := &JWTClaims{}
	cfg := config.Get()
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Auth.JWTSecret), nil
	})

	if err != nil {
		return ctx, errs.NewStandard(connect.CodeUnauthenticated, "解析Token失败")
	}

	if !token.Valid {
		return ctx, errs.NewStandard(connect.CodeUnauthenticated, "Token无效或已过期")
	}

	// 将用户ID存入 context
	ctx = context.WithValue(ctx, constant.UserIDKey, claims.UserID)

	return ctx, nil
}
