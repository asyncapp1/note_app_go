package domain

import "time"

type Note struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	NoteTitle string    `json:"note_title" gorm:"not null"`
	Content   string    `json:"content"`
	IsDone    string    `json:"is_done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteRepository interface {
	Create(note *Note) error
	GetByID(id, userID uint) (*Note, error)
	GetAllByUserID(userID uint) ([]Note, error)
	Update(note *Note, userID uint) error
	Delete(id, userID uint) error
	Query(query string, userID uint) ([]Note, error)
}

type NoteUsecase interface {
	Create(note *Note, user *User) error
	GetByID(id uint, user *User) (*Note, error)
	GetAll(user *User) ([]Note, error)
	Update(note *Note, user *User) error
	Delete(id uint, user *User) error
	Query(query string, user *User) ([]Note, error)
}
