package shortener

import (
	"crypto/rand"
	"math/big"
)

const (
	// Charset for generating short codes (URL-safe characters)
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GenerateCode generates a random short code of specified length
// Uses cryptographically secure random number generator
func GenerateCode(length int) string {
	if length <= 0 {
		length = 7
	}

	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			// Fallback to a less secure but functional method
			result[i] = charset[i%len(charset)]
			continue
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result)
}

// IsValidCode checks if a short code is valid
func IsValidCode(code string) bool {
	if len(code) < 3 || len(code) > 20 {
		return false
	}

	for _, char := range code {
		validChar := false
		for _, c := range charset {
			if char == c {
				validChar = true
				break
			}
		}
		if !validChar {
			return false
		}
	}

	return true
}
