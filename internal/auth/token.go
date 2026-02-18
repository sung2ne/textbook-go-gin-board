package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type TokenService struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	issuer        string
	tokenStore    TokenStore
}

func NewTokenService(
	secretKey string,
	accessExpiry, refreshExpiry time.Duration,
	tokenStore TokenStore,
) *TokenService {
	return &TokenService{
		secretKey:     []byte(secretKey),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		issuer:        "goboard-api",
		tokenStore:    tokenStore,
	}
}

func (s *TokenService) GenerateAccessToken(userID uint, email, username, role string) (string, error) {
	now := time.Now()
	tokenID := generateTokenID()

	claims := CustomClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        tokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *TokenService) GenerateRefreshToken(ctx context.Context, userID uint) (string, error) {
	now := time.Now()
	tokenID := generateTokenID()

	claims := jwt.RegisteredClaims{
		Issuer:    s.issuer,
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ID:        tokenID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	// Redis에 저장
	if s.tokenStore != nil {
		if err := s.tokenStore.StoreRefreshToken(ctx, userID, tokenID, s.refreshExpiry); err != nil {
			return "", err
		}
	}

	return tokenString, nil
}

func (s *TokenService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.secretKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(s.issuer),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, s.handleTokenError(err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *TokenService) ValidateRefreshToken(ctx context.Context, tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return s.secretKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(s.issuer),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return 0, s.handleTokenError(err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return 0, ErrInvalidToken
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		return 0, ErrInvalidToken
	}

	// Redis에서 유효성 확인
	if s.tokenStore != nil {
		valid, err := s.tokenStore.IsRefreshTokenValid(ctx, uint(userID), claims.ID)
		if err != nil {
			return 0, err
		}
		if !valid {
			return 0, ErrInvalidToken
		}
	}

	return uint(userID), nil
}

// RevokeAccessToken은 액세스 토큰을 무효화합니다.
func (s *TokenService) RevokeAccessToken(ctx context.Context, tokenString string) error {
	if s.tokenStore == nil {
		return nil
	}

	token, _ := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return s.secretKey, nil
		},
	)

	if token == nil {
		return ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return ErrInvalidToken
	}

	// 남은 만료 시간 계산
	expiry := time.Until(claims.ExpiresAt.Time)
	if expiry <= 0 {
		return nil
	}

	tokenID := claims.RegisteredClaims.ID
	if tokenID == "" {
		tokenID = HashToken(tokenString)
	}

	return s.tokenStore.AddToBlacklist(ctx, tokenID, expiry)
}

func (s *TokenService) RevokeRefreshToken(ctx context.Context, userID uint) error {
	if s.tokenStore == nil {
		return nil
	}
	return s.tokenStore.DeleteRefreshToken(ctx, userID)
}

// GetAccessExpiry는 액세스 토큰 만료 시간(초)을 반환합니다.
func (s *TokenService) GetAccessExpiry() int64 {
	return int64(s.accessExpiry.Seconds())
}

func (s *TokenService) handleTokenError(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		return fmt.Errorf("%w: malformed token", ErrInvalidToken)
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return fmt.Errorf("%w: invalid signature", ErrInvalidToken)
	case errors.Is(err, jwt.ErrTokenExpired):
		return ErrExpiredToken
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return fmt.Errorf("%w: token not valid yet", ErrInvalidToken)
	default:
		return fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
}

func generateTokenID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// HashToken은 토큰 문자열의 해시를 생성합니다.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:16])
}
