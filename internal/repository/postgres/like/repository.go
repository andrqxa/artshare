package like

import (
	"database/sql"
	"strings"

	modellike "github.com/andrqxa/artshare/internal/model/like"
	"github.com/andrqxa/artshare/internal/repository/postgres"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(connection *postgres.Connection) *Repository {
	return &Repository{db: connection.DB}
}

func (r *Repository) Create(input modellike.CreateInput) error {
	const query = `
		INSERT INTO likes (user_id, artwork_id)
		VALUES ($1, $2)
	`

	if _, err := r.db.Exec(query, input.UserID, input.ArtworkID); err != nil {
		if isUniqueViolation(err) {
			return repoerror.ErrConflict
		}
		return err
	}
	return nil
}

func (r *Repository) Delete(input modellike.DeleteInput) error {
	const query = `
		DELETE FROM likes
		WHERE user_id = $1 AND artwork_id = $2
	`

	result, err := r.db.Exec(query, input.UserID, input.ArtworkID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return repoerror.ErrNotFound
	}
	return nil
}

func (r *Repository) Exists(userID, artworkID string) bool {
	const query = `
		SELECT 1 FROM likes WHERE user_id = $1 AND artwork_id = $2
	`

	var exists int
	if err := r.db.QueryRow(query, userID, artworkID).Scan(&exists); err != nil {
		return false
	}
	return true
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
