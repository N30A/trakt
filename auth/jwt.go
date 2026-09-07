package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/N30A/trakt/config"
	"github.com/N30A/trakt/user"
	"github.com/golang-jwt/jwt/v5"
)

const (
	issuer                = "trakt"
	audience              = "trakt"
	accessTokenExpiration = 15 * time.Minute
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret string
}

func NewJWTService(config config.Config) *JWTService {
	return &JWTService{
		secret: config.JWTSecret,
	}
}

func (s *JWTService) NewAccessToken(id int, role user.UserRole) (string, error) {
	now := time.Now()

	claims := Claims{
		Role:      string(role),
		Issuer:    issuer,
		Subject:   strconv.Itoa(id),
		Audience:  jwt.ClaimStrings{audience},
		ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenExpiration)),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTService) ParseAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	if _, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(s.secret), nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
	); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	return claims, nil
}
