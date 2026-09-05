package auth

import (
	"github.com/N30A/trakt/argon2id"
)

var argon2Params = &argon2id.Params{
	Memory:      argon2id.SecondRecommendedParams.Memory,
	Iterations:  argon2id.SecondRecommendedParams.Iterations,
	Parallelism: argon2id.SecondRecommendedParams.Parallelism,
	SaltLength:  argon2id.SecondRecommendedParams.SaltLength,
	KeyLength:   argon2id.SecondRecommendedParams.KeyLength,
}

func HashPassword(password string) string {
	return argon2id.CreateHash(password, argon2Params)
}

func VerifyPassword(password, hash string) (match bool, err error) {
	match, err = argon2id.ComparePasswordAndHash(password, hash)
	return match, err
}
