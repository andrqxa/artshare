package exchange

import (
	"database/sql"
	"errors"

	modelexchange "github.com/andrqxa/artshare/internal/model/exchange"
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

func (r *Repository) List(filter modelexchange.ListFilter) ([]modelexchange.Summary, pagination.Meta) {
	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	offset := (page - 1) * pageSize

	scope := filter.Scope
	if scope == "" {
		scope = modelexchange.ScopeAll
	}

	status := ""
	if filter.Status != nil {
		status = string(*filter.Status)
	}

	const query = `
		SELECT
			e.id, e.requester_id, e.owner_id, e.status, e.message,
			e.created_at, e.updated_at,
			a.id, a.title, a.preview_image_url,
			COUNT(*) OVER() AS total_items
		FROM exchange_requests e
		JOIN artworks a ON a.id = e.artwork_id
		WHERE ($1 = '' OR
			($2 = 'inbox' AND e.owner_id::text = $1) OR
			($2 = 'outbox' AND e.requester_id::text = $1) OR
			($2 = 'all' AND (e.owner_id::text = $1 OR e.requester_id::text = $1)))
		  AND ($3 = '' OR e.status = $3)
		ORDER BY e.created_at DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := r.db.Query(query, filter.UserID, string(scope), status, pageSize, offset)
	if err != nil {
		return nil, paginationMeta(page, pageSize, 0)
	}
	defer rows.Close()

	totalItems := 0
	items := make([]modelexchange.Summary, 0)
	for rows.Next() {
		var (
			s            modelexchange.Summary
			message      sql.NullString
			previewImage sql.NullString
		)
		if err := rows.Scan(
			&s.ID, &s.RequesterID, &s.OwnerID, &s.Status, &message,
			&s.CreatedAt, &s.UpdatedAt,
			&s.Artwork.ID, &s.Artwork.Title, &previewImage,
			&totalItems,
		); err != nil {
			return nil, paginationMeta(page, pageSize, 0)
		}
		s.Artwork.OwnerID = s.OwnerID
		s.Artwork.PreviewImageURL = nullString(previewImage)
		items = append(items, s)
	}
	if err := rows.Err(); err != nil {
		return nil, paginationMeta(page, pageSize, 0)
	}
	return items, paginationMeta(page, pageSize, totalItems)
}

func (r *Repository) Create(input modelexchange.CreateInput) (modelexchange.Detail, error) {
	const insertQuery = `
		INSERT INTO exchange_requests (artwork_id, requester_id, owner_id, status, message)
		SELECT $1, $2, ar.user_id, 'pending', $3
		FROM artworks a
		JOIN artists ar ON ar.id = a.artist_id
		WHERE a.id = $1
		RETURNING id
	`

	var id string
	err := r.db.QueryRow(insertQuery, input.ArtworkID, input.RequesterID, input.Message).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return modelexchange.Detail{}, repoerror.ErrNotFound
	}
	if err != nil {
		return modelexchange.Detail{}, err
	}
	return r.FindByID(id)
}

func (r *Repository) FindByID(id string) (modelexchange.Detail, error) {
	const query = `
		SELECT
			e.id, e.requester_id, e.owner_id, e.status, e.message,
			e.created_at, e.updated_at,
			a.id, a.title, a.preview_image_url
		FROM exchange_requests e
		JOIN artworks a ON a.id = e.artwork_id
		WHERE e.id = $1
	`

	var (
		detail        modelexchange.Detail
		message       sql.NullString
		previewImage  sql.NullString
	)
	err := r.db.QueryRow(query, id).Scan(
		&detail.ID, &detail.RequesterID, &detail.OwnerID, &detail.Status, &message,
		&detail.CreatedAt, &detail.UpdatedAt,
		&detail.Artwork.ID, &detail.Artwork.Title, &previewImage,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return modelexchange.Detail{}, repoerror.ErrNotFound
	}
	if err != nil {
		return modelexchange.Detail{}, err
	}
	detail.Message = nullString(message)
	detail.Artwork.OwnerID = detail.OwnerID
	detail.Artwork.PreviewImageURL = nullString(previewImage)
	return detail, nil
}

func (r *Repository) UpdateStatus(input modelexchange.UpdateStatusInput) (modelexchange.Detail, error) {
	const query = `
		UPDATE exchange_requests
		SET status = $2,
		    message = COALESCE($3, message),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id
	`

	var id string
	err := r.db.QueryRow(query, input.ID, string(input.Status), input.Message).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return modelexchange.Detail{}, repoerror.ErrNotFound
	}
	if err != nil {
		return modelexchange.Detail{}, err
	}
	return r.FindByID(id)
}

func nullString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	v := value.String
	return &v
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
