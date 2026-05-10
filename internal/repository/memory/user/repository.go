package user

import (
	"fmt"
	"sync"
	"time"

	modeluser "github.com/andrqxa/artshare/internal/model/user"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Repository struct {
	mu          sync.RWMutex
	nextID      int
	byEmail     map[string]record
	currentUser modeluser.CurrentUser
}

type record struct {
	User     modeluser.CurrentUser
	Password string
}

func NewRepository() *Repository {
	currentUser := modeluser.CurrentUser{
		ID:        "00000000-0000-0000-0000-000000000001",
		Email:     "viewer@artshare.local",
		Role:      modeluser.RoleViewer,
		CreatedAt: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	return &Repository{
		nextID:      2,
		currentUser: currentUser,
		byEmail: map[string]record{
			currentUser.Email: {
				User:     currentUser,
				Password: "password",
			},
		},
	}
}

func (r *Repository) Create(input modeluser.RegisterInput) (modeluser.CurrentUser, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byEmail[input.Email]; exists {
		return modeluser.CurrentUser{}, repoerror.ErrConflict
	}

	user := modeluser.CurrentUser{
		ID:        nextUUID(r.nextID),
		Email:     input.Email,
		Role:      input.Role,
		CreatedAt: time.Now().UTC(),
	}
	r.nextID++
	r.byEmail[user.Email] = record{
		User:     user,
		Password: input.Password,
	}
	r.currentUser = user

	return user, nil
}

func (r *Repository) FindByEmail(email string) (modeluser.CurrentUser, string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	record, exists := r.byEmail[email]
	if !exists {
		return modeluser.CurrentUser{}, "", repoerror.ErrNotFound
	}
	return record.User, record.Password, nil
}

func (r *Repository) GetCurrent() modeluser.CurrentUser {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.currentUser
}

func (r *Repository) SetCurrent(user modeluser.CurrentUser) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.currentUser = user
}

func (r *Repository) UpdateCurrent(input modeluser.UpdateCurrentUserInput) modeluser.CurrentUser {
	r.mu.Lock()
	defer r.mu.Unlock()

	if input.Email != nil {
		currentRecord := r.byEmail[r.currentUser.Email]
		delete(r.byEmail, r.currentUser.Email)
		r.currentUser.Email = *input.Email
		currentRecord.User = r.currentUser
		r.byEmail[r.currentUser.Email] = currentRecord
	}

	return r.currentUser
}

func nextUUID(id int) string {
	return fmt.Sprintf("00000000-0000-0000-0000-%012d", id)
}
