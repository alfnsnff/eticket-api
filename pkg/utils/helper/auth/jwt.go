package auth // Place this in a shared utility package

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5" // Import the JWT library
)

// Define your JWT signing key. Keep this secret and secure!
// In a real application, load this from environment variables or a secure configuration.
var jwtSecretKey = []byte("your_super_secret_jwt_key") // CHANGE THIS IN PRODUCTION

// Claims struct represents the custom claims you want to include in the JWT payload.
// Standard claims are embedded using jwt.RegisteredClaims.
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	// Add other custom claims here if needed (e.g., Role, Permissions)
	jwt.RegisteredClaims // Embed standard JWT claims (Issuer, Subject, Audience, ExpiresAt, etc.)
}

// GenerateAccessToken generates a new JWT access token.
func GenerateAccessToken(userID uint, username string) (string, error) {
	// Define the expiry time for the access token
	// Access tokens should generally be short-lived for security
	expirationTime := time.Now().Add(15 * time.Minute) // Example: 15 minutes expiry

	// Create the claims
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // Set expiry time
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // Set issued time
			// Issuer:    "your_app_name", // Optional: Set issuer
			// Subject:   fmt.Sprintf("%d", userID), // Optional: Set subject (often user ID)
		},
	}

	// Create the token using your claims and the signing method (e.g., HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a new JWT refresh token.
// Refresh tokens are typically longer-lived and used to obtain new access tokens
// without requiring the user to log in again.
// You might store refresh tokens in the database and invalidate them on logout.
func GenerateRefreshToken(userID uint) (string, error) {
	// Define the expiry time for the refresh token
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // Example: 7 days expiry

	// Create the claims (often simpler claims than access tokens)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			// Add a unique identifier (JTI) to the refresh token claims
			// This helps in managing refresh token validity (e.g., revoking)
			// Jti: uuid.New().String(), // Requires importing github.com/google/uuid
		},
	}

	// Create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token string and returns the claims if valid.
func ValidateToken(tokenString string) (*Claims, error) {
	// Parse the token string
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for validation
		return jwtSecretKey, nil
	})

	// Handle parsing errors
	if err != nil {
		// Check if the error is a validation error (e.g., expired token)
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		// Handle other parsing errors (e.g., malformed token)
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract and return the claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("failed to extract claims")
	}

	// Optional: Check if the token is expired (jwt.ParseWithClaims usually handles this if ExpiresAt is set)
	// if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
	// 	return nil, errors.New("token expired")
	// }

	return claims, nil
}

// --- Example Usage ---
/*
// In a login use case after successful password check:
// accessToken, err := utils.GenerateAccessToken(user.ID, user.Username)
// refreshToken, err := utils.GenerateRefreshToken(user.ID)
// // Return tokens to the client

// In an authentication middleware:
// tokenString := // Get token from Authorization header
// claims, err := utils.ValidateToken(tokenString)
// if err != nil { return unauthorized error }
// // Token is valid, claims contain user_id and username
// // You can now use claims.UserID or claims.Username
*/
