package usecase

import "notes-app/internal/domain"

type noteUsecase struct {
	noteRepo domain.NoteRepository
}

func NewNoteUsecase(repo domain.NoteRepository) domain.NoteUsecase {
	return &noteUsecase{
		noteRepo: repo,
	}
}

func (u *noteUsecase) Create(note *domain.Note, user *domain.User) error {
	note.UserID = user.ID
	return u.noteRepo.Create(note)
}

func (u *noteUsecase) GetByID(id uint, user *domain.User) (*domain.Note, error) {
	return u.noteRepo.GetByID(id, user.ID)
}

func (u *noteUsecase) GetAll(user *domain.User) ([]domain.Note, error) {
	return u.noteRepo.GetAllByUserID(user.ID)
}

func (u *noteUsecase) Update(note *domain.Note, user *domain.User) error {
	return u.noteRepo.Update(note, user.ID)
}

func (u *noteUsecase) Delete(id uint, user *domain.User) error {
	return u.noteRepo.Delete(id, user.ID)
}

func (u *noteUsecase) Query(query string, user *domain.User) ([]domain.Note, error) {
	return u.noteRepo.Query(query, user.ID)
}
