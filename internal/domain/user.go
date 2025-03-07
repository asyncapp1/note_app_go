package domain

import (
	"notes-app/pkg/types"
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"` // "-" means don't show in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByUsername(username string) (*User, error)
	GetByID(id uint) (*User, error)
}

type UserUsecase interface {
	Register(user *User) error
	Login(username, password string) (*types.TokenPair, error)
	RefreshToken(refreshToken string) (*types.TokenPair, error)
	ValidateAccessToken(token string) (*User, error)
}
