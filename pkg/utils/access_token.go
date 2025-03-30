package utils

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/pkg/logger"
)

var (
	ErrorInvalidToken = errors.New("invalid token")
	AuthorizationKey  = "authorization"
	RefreshTokenKey   = "refresh_token"
)

func GetAuthorizationToken(r *http.Request) string {
	tokenWithBearer := r.Header.Get(AuthorizationKey)
	onlyToken := strings.Split(tokenWithBearer, "Bearer ")
	if len(onlyToken) < 2 {
		return ""
	}
	return onlyToken[1]
}

type AccessTokenClaims struct {
	UserID uuid.UUID      `json:"user_id"`
	Email  string         `json:"email"`
	Rol    model.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func VerifyAccessToken(tokenString string, key []byte) (*AccessTokenClaims, error) {
	claims := &AccessTokenClaims{}
	logs := logger.GetGlobalLogger()
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		logs.Error(err)
		return nil, ErrorInvalidToken
	}
	if !token.Valid {
		return nil, ErrorInvalidToken
	}
	return claims, nil
}

func GenerateAccessToken(user *model.User, accessExpirationTime time.Duration, key []byte) (string, error) {
	claims := AccessTokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Rol:    user.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpirationTime)),
		},
	}

	logs := logger.GetGlobalLogger()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	return tokenString, nil
}

func GetClaimsFromContext(ctx context.Context) *AccessTokenClaims {
	claims, ok := ctx.Value(AuthorizationKey).(*AccessTokenClaims)
	if !ok {
		return nil
	}
	return claims
}

func NewCookieRefreshToken(token string, maxAge time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     RefreshTokenKey,
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   int(maxAge.Seconds()),
	}
}

func GetRefreshTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(RefreshTokenKey)
	if err != nil {
		return ""
	}
	return cookie.Value
}
