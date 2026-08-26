package jwt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"mrb-service/internal/config"
)

type RefreshToken struct {
	Token     string
	ExpiresAt time.Time
}

const (
	Access  string = "access"
	Refresh string = "refresh"
)

type Manager interface {
	GenerateRefreshToken() (string, error)
	CreateAccessToken(userID, role string, permissions []string) (string, error)
	ValidateAccessToken(tokenString string) (*claims, error)
	HashToken(refreshToken string) string
	RefreshTokenExpiresAt() time.Time
}

type claims struct {
	UserID      string   `json:"uid"`
	Role        string   `json:"rol"`
	Permissions []string `json:"permissions"`
	TokenType   string   `json:"typ"`

	jwt.RegisteredClaims
}

func newClaims(
	userID, role, tokenType, issuer string,
	permissions []string,
	accessMinutes time.Duration,
) *claims {
	return &claims{
		UserID:      userID,
		Role:        role,
		Permissions: permissions,
		TokenType:   tokenType,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessMinutes * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   userID,
		},
	}
}

type manager struct {
	config      config.JWT
	accessName  string
	refreshName string
}

func New(cfg config.JWT) *manager {
	return &manager{
		config:      cfg,
		accessName:  Access,
		refreshName: Refresh,
	}
}

func (m *manager) CreateAccessToken(userID, role string, permissions []string) (string, error) {
	claim := newClaims(
		userID,
		role,
		m.accessName,
		m.config.GetIssuer(),
		permissions,
		m.config.GetAccessExpires(),
	)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	accessToken, err := token.SignedString(m.config.GetPrivateKey())
	if err != nil {
		return accessToken, err
	}

	return accessToken, nil
}

func (m *manager) ValidateAccessToken(tokenString string) (*claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodRS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.config.GetPublicKey(), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*claims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if claims.TokenType != Access {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

func (_ *manager) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 48)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (_ *manager) HashToken(refreshToken string) string {
	hash := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(hash[:])
}

func (m *manager) RefreshTokenExpiresAt() time.Time {
	return time.Now().Add(m.config.GetRefreshExpires())
}
