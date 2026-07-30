package auth

import (
	"golang.org/x/crypto/argon2"
	"github.com/ieraHQ/Vakalat/backend/api/config"
)

// HashPassword hashes a password using Argon2id.
func HashPassword(password string, cfg *config.Config) (string, error) {
	salt := []byte(cfg.JWT.Secret) // Use JWT secret as salt for simplicity
	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)
	return string(hash), nil
}

// VerifyPassword verifies a password against a hash.
func VerifyPassword(password, hash string, cfg *config.Config) bool {
	salt := []byte(cfg.JWT.Secret) // Use JWT secret as salt for simplicity
	newHash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)
	return string(newHash) == hash
}