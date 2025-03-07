package repository

import (
	"errors"
	"notes-app/internal/domain"

	"gorm.io/gorm"
)

type noteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) domain.NoteRepository {
	return &noteRepository{db}
}

func (r *noteRepository) Create(note *domain.Note) error {
	return r.db.Create(note).Error
}

func (r *noteRepository) GetByID(id, userID uint) (*domain.Note, error) {
	var note domain.Note
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&note).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("note not found or unauthorized")
		}
		return nil, err
	}
	return &note, nil
}

func (r *noteRepository) GetAllByUserID(userID uint) ([]domain.Note, error) {
	var notes []domain.Note
	err := r.db.Where("user_id = ?", userID).Find(&notes).Error
	return notes, err
}

func (r *noteRepository) Update(note *domain.Note, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", note.ID, userID).Updates(note)
	if result.RowsAffected == 0 {
		return errors.New("note not found or unauthorized")
	}
	return result.Error
}

func (r *noteRepository) Delete(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Note{})
	if result.RowsAffected == 0 {
		return errors.New("note not found or unauthorized")
	}
	return result.Error
}

func (r *noteRepository) Query(query string, userID uint) ([]domain.Note, error) {
	var notes []domain.Note
	err := r.db.Where("user_id = ? AND (note_title ILIKE ? OR content ILIKE ?)",
		userID,
		"%"+query+"%",
		"%"+query+"%").
		Find(&notes).Error
	return notes, err
}
