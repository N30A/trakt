package auth

import (
	"context"
	"errors"

	"github.com/N30A/trakt/argon2id"
	"github.com/N30A/trakt/user"
)

var (
	argon2Params = argon2id.Params{
		Memory:      argon2id.SecondRecommendedParams.Memory,
		Iterations:  argon2id.SecondRecommendedParams.Iterations,
		Parallelism: argon2id.SecondRecommendedParams.Parallelism,
		SaltLength:  argon2id.SecondRecommendedParams.SaltLength,
		KeyLength:   argon2id.SecondRecommendedParams.KeyLength,
	}
)

type AuthService struct {
	userRepo   *user.UserRepo
	jwtService *JWTService
}

func NewAuthService(userRepo *user.UserRepo, jwtService *JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) CreateInitialUser(ctx context.Context, email, password string) (user.User, error) {
	passwordHash := argon2id.CreateHash(password, argon2Params)

	newUser, err := s.userRepo.CreateInitialUser(ctx, email, passwordHash)
	if err != nil {
		if errors.Is(err, user.ErrConflict) {
			return user.User{}, ErrInitialUserExists
		}
		return user.User{}, err
	}

	return newUser, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	usr, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	match, err := argon2id.ComparePasswordAndHash(password, usr.PasswordHash)
	if err != nil {
		return "", err
	}

	if !match {
		return "", ErrInvalidCredentials
	}

	return s.jwtService.NewAccessToken(usr.ID, usr.Role)
}
