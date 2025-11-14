// Copyright 2025 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	// BcryptCost defines the cost factor for bcrypt hashing
	// Cost 12 provides a good balance between security and performance
	BcryptCost = 12

	// DefaultPasswordResetTokenLength defines the length of password reset tokens
	DefaultPasswordResetTokenLength = 32
)

// HashPassword hashes a password using bcrypt with the configured cost
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// ValidatePassword compares a password with its hash using bcrypt
func ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordComplexity validates password complexity based on configuration
func ValidatePasswordComplexity(password string, config PasswordAuthConfig) error {
	if len(password) < config.Policy.MinLength {
		return fmt.Errorf("password must be at least %d characters long", config.Policy.MinLength)
	}

	if config.Policy.RequireUppercase && !containsUppercase(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if config.Policy.RequireLowercase && !containsLowercase(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if config.Policy.RequireNumbers && !containsNumbers(password) {
		return fmt.Errorf("password must contain at least one number")
	}

	if config.Policy.RequireSymbols && !containsSymbols(password) {
		return fmt.Errorf("password must contain at least one symbol")
	}

	return nil
}

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken() (string, error) {
	return GenerateSecureTokenWithLength(DefaultPasswordResetTokenLength)
}

// GenerateSecureTokenWithLength generates a secure token with specified length
func GenerateSecureTokenWithLength(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// containsUppercase checks if the string contains uppercase letters
func containsUppercase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

// containsLowercase checks if the string contains lowercase letters
func containsLowercase(s string) bool {
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return true
		}
	}
	return false
}

// containsNumbers checks if the string contains numbers
func containsNumbers(s string) bool {
	re := regexp.MustCompile(`\d`)
	return re.MatchString(s)
}

// containsSymbols checks if the string contains symbols
func containsSymbols(s string) bool {
	// Check for common symbols
	symbols := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	return strings.ContainsAny(s, symbols)
}
