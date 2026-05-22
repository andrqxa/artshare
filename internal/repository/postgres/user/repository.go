package user

import (
	"database/sql"
	"errors"
	"strings"
	"sync"

	modeluser "github.com/andrqxa/artshare/internal/model/user"
	"github.com/andrqxa/artshare/internal/repository/postgres"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Repository struct {
	db *sql.DB

	mu          sync.RWMutex
	currentUser modeluser.CurrentUser
}

func NewRepository(connection *postgres.Connection) *Repository {
	return &Repository{db: connection.DB}
}

func (r *Repository) Create(input modeluser.RegisterInput) (modeluser.CurrentUser, error) {
	const query = `
		INSERT INTO users (email, password, role)
		VALUES ($1, $2, $3)
		RETURNING id, email, role, created_at
	`

	user, err := scanCurrentUser(r.db.QueryRow(query, input.Email, input.Password, string(input.Role)))
	if err != nil {
		if isUniqueViolation(err) {
			return modeluser.CurrentUser{}, repoerror.ErrConflict
		}
		return modeluser.CurrentUser{}, err
	}

	r.SetCurrent(user)
	return user, nil
}

func (r *Repository) FindByEmail(email string) (modeluser.CurrentUser, string, error) {
	const query = `
		SELECT id, email, role, created_at, password
		FROM users
		WHERE email = $1
	`

	var (
		user     modeluser.CurrentUser
		role     string
		password string
	)
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &role, &user.CreatedAt, &password)
	if errors.Is(err, sql.ErrNoRows) {
		return modeluser.CurrentUser{}, "", repoerror.ErrNotFound
	}
	if err != nil {
		return modeluser.CurrentUser{}, "", err
	}
	user.Role = modeluser.Role(role)
	return user, password, nil
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

	if input.Email != nil && r.currentUser.ID != "" {
		const query = `UPDATE users SET email = $2, updated_at = NOW() WHERE id = $1`
		if _, err := r.db.Exec(query, r.currentUser.ID, *input.Email); err == nil {
			r.currentUser.Email = *input.Email
		}
	}

	return r.currentUser
}

func scanCurrentUser(row *sql.Row) (modeluser.CurrentUser, error) {
	var (
		user modeluser.CurrentUser
		role string
	)
	if err := row.Scan(&user.ID, &user.Email, &role, &user.CreatedAt); err != nil {
		return modeluser.CurrentUser{}, err
	}
	user.Role = modeluser.Role(role)
	return user, nil
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
