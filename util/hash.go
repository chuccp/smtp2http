package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// MinPasswordLength is the shortest password the API accepts. The instance is normally
// reachable from the network and the only other rule is "not empty" — both the forms and
// the handlers used to accept a single character.
const MinPasswordLength = 8

// ValidatePassword rejects a password shorter than MinPasswordLength. Counted in runes,
// so a password of CJK characters is measured the way a person would count it.
func ValidatePassword(password string) error {
	if len([]rune(password)) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters", MinPasswordLength)
	}
	return nil
}

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plaintext password with a bcrypt hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
