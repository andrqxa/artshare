package artist

import (
	"database/sql"
	"errors"
	"strings"

	modelartist "github.com/andrqxa/artshare/internal/model/artist"
	"github.com/andrqxa/artshare/internal/model/pagination"
	"github.com/andrqxa/artshare/internal/repository/postgres"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(connection *postgres.Connection) *Repository {
	return &Repository{db: connection.DB}
}

func (r *Repository) Create(input modelartist.CreateInput) (modelartist.Detail, error) {
	const query = `
		INSERT INTO artists (user_id, display_name, bio, avatar_url)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, display_name, bio, avatar_url, created_at, updated_at
	`

	artist, err := scanDetail(r.db.QueryRow(query, input.UserID, input.DisplayName, input.Bio, input.AvatarURL))
	if err != nil {
		if isUniqueViolation(err) {
			return modelartist.Detail{}, repoerror.ErrConflict
		}
		return modelartist.Detail{}, err
	}

	return artist, nil
}

func (r *Repository) FindByUserID(userID string) (modelartist.Detail, error) {
	const query = `
		SELECT id, user_id, display_name, bio, avatar_url, created_at, updated_at
		FROM artists
		WHERE user_id = $1
	`

	return detailOrRepositoryError(scanDetail(r.db.QueryRow(query, userID)))
}

func (r *Repository) FindByID(id string) (modelartist.Detail, error) {
	const query = `
		SELECT id, user_id, display_name, bio, avatar_url, created_at, updated_at
		FROM artists
		WHERE id = $1
	`

	return detailOrRepositoryError(scanDetail(r.db.QueryRow(query, id)))
}

func (r *Repository) UpdateByUserID(userID string, input modelartist.UpdateInput) (modelartist.Detail, error) {
	const query = `
		UPDATE artists
		SET
			display_name = COALESCE($2, display_name),
			bio = COALESCE($3, bio),
			avatar_url = COALESCE($4, avatar_url),
			updated_at = NOW()
		WHERE user_id = $1
		RETURNING id, user_id, display_name, bio, avatar_url, created_at, updated_at
	`

	return detailOrRepositoryError(scanDetail(r.db.QueryRow(
		query,
		userID,
		input.DisplayName,
		input.Bio,
		input.AvatarURL,
	)))
}

func (r *Repository) List(filter modelartist.ListFilter) ([]modelartist.Summary, pagination.Meta) {
	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	offset := (page - 1) * pageSize
	query := strings.TrimSpace(valueOrEmpty(filter.Query))

	const sqlQuery = `
		SELECT
			a.id,
			a.display_name,
			a.avatar_url,
			a.bio,
			COUNT(aw.id) AS artwork_count,
			COUNT(*) OVER() AS total_items
		FROM artists a
		LEFT JOIN artworks aw ON aw.artist_id = a.id
		WHERE $1 = '' OR a.display_name ILIKE '%' || $1 || '%'
		GROUP BY a.id, a.display_name, a.avatar_url, a.bio
		ORDER BY a.display_name
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(sqlQuery, query, pageSize, offset)
	if err != nil {
		return nil, paginationMeta(page, pageSize, 0)
	}
	defer rows.Close()

	var totalItems int
	items := make([]modelartist.Summary, 0)
	for rows.Next() {
		var artist modelartist.Summary
		if err := rows.Scan(
			&artist.ID,
			&artist.DisplayName,
			&artist.AvatarURL,
			&artist.Bio,
			&artist.ArtworkCount,
			&totalItems,
		); err != nil {
			return nil, paginationMeta(page, pageSize, 0)
		}
		items = append(items, artist)
	}
	if err := rows.Err(); err != nil {
		return nil, paginationMeta(page, pageSize, 0)
	}

	return items, paginationMeta(page, pageSize, totalItems)
}

func scanDetail(row *sql.Row) (modelartist.Detail, error) {
	var artist modelartist.Detail
	err := row.Scan(
		&artist.ID,
		&artist.UserID,
		&artist.DisplayName,
		&artist.Bio,
		&artist.AvatarURL,
		&artist.CreatedAt,
		&artist.UpdatedAt,
	)
	return artist, err
}

func detailOrRepositoryError(artist modelartist.Detail, err error) (modelartist.Detail, error) {
	if errors.Is(err, sql.ErrNoRows) {
		return modelartist.Detail{}, repoerror.ErrNotFound
	}
	return artist, err
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func normalizePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func paginationMeta(page, pageSize, totalItems int) pagination.Meta {
	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}
	return pagination.Meta{Page: page, PageSize: pageSize, TotalItems: totalItems, TotalPages: totalPages}
}
