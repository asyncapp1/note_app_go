package usecase

import (
	"errors"
	"log"
	"notes-app/internal/domain"
	"notes-app/pkg/auth"
	"notes-app/pkg/types"

	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: repo,
	}
}

func (u *userUsecase) Register(user *domain.User) error {
	// Log original password (remove in production)
	log.Printf("Registering user - Password before hash length: %d", len(user.Password))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Log hashed password (remove in production)
	log.Printf("Generated hash length: %d", len(hashedPassword))

	user.Password = string(hashedPassword)
	return u.userRepo.Create(user)
}

func (u *userUsecase) Login(username, password string) (*types.TokenPair, error) {
	log.Printf("Login usecase - Username: %s, Password length: %d",
		username, len(password))

	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return nil, errors.New("invalid credentials")
	}

	log.Printf("Found user with ID: %d, Stored password hash length: %d", user.ID, len(user.Password))

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password comparison failed: %v", err)
		return nil, errors.New("invalid credentials")
	}

	authToken, err := auth.GenerateTokenPair(user.ID)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		return nil, err
	}

	return &types.TokenPair{
		AccessToken:  authToken.AccessToken,
		RefreshToken: authToken.RefreshToken,
	}, nil
}

func (u *userUsecase) RefreshToken(refreshToken string) (*types.TokenPair, error) {
	claims, err := auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	authToken, err := auth.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, err
	}

	return &types.TokenPair{
		AccessToken:  authToken.AccessToken,
		RefreshToken: authToken.RefreshToken,
	}, nil
}

func (u *userUsecase) ValidateAccessToken(token string) (*domain.User, error) {
	claims, err := auth.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}

	return u.userRepo.GetByID(claims.UserID)
}
