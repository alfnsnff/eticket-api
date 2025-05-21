package utils // Place this in a shared utility package

import (
	"fmt"

	"golang.org/x/crypto/bcrypt" // Import bcrypt library
)

// HashPassword securely hashes a plain text password using bcrypt.
func HashPassword(password string) (string, error) {
	// GenerateFromPassword takes the password and the cost factor.
	// bcrypt.DefaultCost is a good default, but you can adjust it based on your needs and server resources.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a plain text password with a hashed password.
// It returns true if they match, false otherwise.
func CheckPasswordHash(password, hash string) bool {
	// CompareHashAndPassword compares the hashed password with its possible plain-text equivalent.
	// It returns nil on success, or an error on failure (including incorrect password).
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // Return true if err is nil (passwords match)
}

// --- Example Usage ---
/*
// In a registration use case:
// hashedPassword, err := utils.HashPassword(plainPassword)
// if err != nil { handle error }
// user.Password = hashedPassword
// // Save user to database

// In a login use case:
// user, err := userRepository.GetUserByUsername(username)
// if err != nil { handle error }
// if user == nil { return errors.New("user not found") }
// if !utils.CheckPasswordHash(plainPassword, user.Password) {
//     return errors.New("invalid password")
// }
// // Password is correct, proceed with token generation
*/
