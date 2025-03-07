package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Type   string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

func GenerateTokenPair(userID uint) (*TokenPair, error) {
	// Generate Access Token (short-lived)
	accessClaims := JWTClaims{
		UserID: userID,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token (long-lived)
	refreshClaims := JWTClaims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	return validateToken(tokenString, "access", os.Getenv("JWT_SECRET"))
}

func ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	return validateToken(tokenString, "refresh", os.Getenv("JWT_REFRESH_SECRET"))
}

func validateToken(tokenString, tokenType, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		if claims.Type != tokenType {
			return nil, errors.New("invalid token type")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
