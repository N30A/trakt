package validate

import (
	"errors"
	"net/mail"
	"strings"
)

func Email(value string) (string, error) {
	value = strings.TrimSpace(value)

	parsed, err := mail.ParseAddress(value)
	if err != nil || parsed.Address != value {
		return "", errors.New("invalid email format")
	}

	return parsed.Address, nil
}
