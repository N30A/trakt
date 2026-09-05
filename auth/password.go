package auth

// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#argon2id
// https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
// https://github.com/p-h-c/phc-winner-argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	memory      uint32 = 46 * 1024
	iterations  uint32 = 1
	parallelism uint8  = 1
	keyLength   uint32 = 32
	saltLength  uint32 = 16
)

type verifyParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

func generateSalt(saltLength uint32) []byte {
	b := make([]byte, saltLength)
	rand.Read(b)
	return b
}

func HashPassword(password string) string {
	salt := generateSalt(saltLength)
	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	// Base64 without padding
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// PHC format: $argon2id$v=<version>$m=<memory>,t=<iterations>,p=<parallelism>$<b64salt>$<b64hash>
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, iterations, parallelism, b64Salt, b64Hash)
}

func decodeHash(encodedHash string) (salt, hash []byte, params *verifyParams) {
	values := strings.Split(encodedHash, "$")
	if len(values) != 6 {
		return nil, nil, nil
	}

	var version int
	_, err := fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil
	}

	if version != argon2.Version {
		return nil, nil, nil
	}

	params = &verifyParams{}

	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism)
	if err != nil {
		return nil, nil, nil
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return nil, nil, nil
	}

	hash, err = base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return nil, nil, nil
	}
	params.keyLength = uint32(len(hash))

	return salt, hash, params
}

func VerifyPassword(password, encodedHash string) bool {
	salt, hash, params := decodeHash(encodedHash)
	if salt == nil || hash == nil || params == nil {
		return false
	}

	otherHash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)
	return subtle.ConstantTimeCompare(hash, otherHash) == 1
}
