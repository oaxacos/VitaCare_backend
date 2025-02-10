package token

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/oaxacos/vitacare/internal/config"
	"github.com/oaxacos/vitacare/internal/domain/model"
	"github.com/oaxacos/vitacare/internal/domain/repository"
	"github.com/oaxacos/vitacare/pkg/logger"
)

type AccessTokenClaims struct {
	UserID uuid.UUID      `json:"user_id"`
	Email  string         `json:"email"`
	Rol    model.UserRole `json:"role"`
	jwt.RegisteredClaims
}

var (
	ErrInvalidToken = errors.New("invalid token")
)

type TokenService struct {
	accessTokenKey        []byte
	refreshTokenKey       []byte
	accessExpirationTime  time.Duration
	refreshExpirationTime time.Duration
	repo                  repository.RefreshTokenRepository
}

func NewTokenService(conf *config.Config, repo repository.RefreshTokenRepository) *TokenService {
	return &TokenService{
		accessTokenKey:        []byte(conf.Token.PrivateKeyAccessToken),
		refreshTokenKey:       []byte(conf.Token.PrivateKeyRefreshToken),
		accessExpirationTime:  time.Duration(conf.Token.AccessTimeExpiration) * time.Minute,
		refreshExpirationTime: time.Duration(conf.Token.RefreshTimeExpiration) * time.Hour,
		repo:                  repo,
	}
}

func (t *TokenService) GenerateAccessToken(ctx context.Context, user *model.User) (string, error) {
	claims := AccessTokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		Rol:    user.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.accessExpirationTime)),
		},
	}

	logs := logger.GetContextLogger(ctx)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.accessTokenKey)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	return tokenString, nil
}

func (t *TokenService) GenerateRefreshToken(ctx context.Context, user *model.User) (string, error) {
	logs := logger.GetContextLogger(ctx)

	tokenString, err := t.generateRandomToken()
	if err != nil {
		logs.Error(err)
		return "", err
	}

	refreshToken := model.NewRefreshToken(tokenString, user.ID, t.refreshExpirationTime)
	err = t.repo.Save(ctx, refreshToken)
	if err != nil {
		logs.Error(err)
		return "", err
	}
	return tokenString, nil
}

func (t *TokenService) VerifyAccessToken(ctx context.Context, token string) (*AccessTokenClaims, error) {
	return t.validateToken(ctx, token, string(t.accessTokenKey))
}

func (t *TokenService) VerifyRefreshToken(ctx context.Context, token string) error {
	logs := logger.GetContextLogger(ctx)
	if token == "" {
		return ErrInvalidToken
	}
	tokenInDB, err := t.repo.GetByToken(ctx, token)
	if err != nil {
		logs.Error(err)
		return err
	}

	if t.isExpired(tokenInDB) {
		return ErrInvalidToken
	}
	return nil
}
func (t *TokenService) isExpired(token *model.RefreshToken) bool {
	return token.ExpiredAt.Before(time.Now())
}

func (t *TokenService) validateToken(ctx context.Context, tokenString, secret string) (*AccessTokenClaims, error) {
	logs := logger.GetContextLogger(ctx)
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logs.Error("unexpected signing method")
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		logs.Error(err)
		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		logs.Error("invalid token")
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (t *TokenService) generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (t *TokenService) RenewTokens(ctx context.Context, user *model.User) (string, string, error) {
	log := logger.GetContextLogger(ctx)
	beforeToken, err := t.repo.GetByUserID(ctx, user.ID)
	if err != nil {
		log.Error(err)
		return "", "", err
	}
	if beforeToken != nil {
		err = t.repo.Delete(ctx, beforeToken.ID)
		if err != nil {
			return "", "", err
		}
	}
	accessToken, err := t.GenerateAccessToken(ctx, user)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := t.GenerateRefreshToken(ctx, user)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
