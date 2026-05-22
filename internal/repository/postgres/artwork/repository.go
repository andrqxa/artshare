package artwork

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	modelartwork "github.com/andrqxa/artshare/internal/model/artwork"
	"github.com/andrqxa/artshare/internal/model/pagination"
	"github.com/andrqxa/artshare/internal/repository/postgres"
	"github.com/andrqxa/artshare/internal/repository/repoerror"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(connection *postgres.Connection) *Repository {
	return &Repository{db: connection.DB}
}

func (r *Repository) List(filter modelartwork.ListFilter) ([]modelartwork.Summary, pagination.Meta) {
	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	offset := (page - 1) * pageSize

	query := strings.TrimSpace(valueOrEmpty(filter.Query))
	artistID := valueOrEmpty(filter.ArtistID)
	status := ""
	if filter.Status != nil {
		status = string(*filter.Status)
	}
	order := "DESC"
	if filter.Sort == "oldest" {
		order = "ASC"
	}

	sqlQuery := `
		SELECT
			a.id, a.title, a.status, a.preview_image_url,
			a.availability_for_exchange, a.created_at, a.published_at,
			ar.id, ar.display_name, ar.avatar_url,
			(SELECT COUNT(*) FROM likes l WHERE l.artwork_id = a.id) AS like_count,
			COUNT(*) OVER() AS total_items
		FROM artworks a
		JOIN artists ar ON ar.id = a.artist_id
		WHERE ($1 = '' OR a.title ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR a.artist_id::text = $2)
		  AND ($3 = '' OR a.status = $3)
		ORDER BY a.created_at ` + order + `
		LIMIT $4 OFFSET $5
	`

	rows, err := r.db.Query(sqlQuery, query, artistID, status, pageSize, offset)
	if err != nil {
		return nil, paginationMeta(page, pageSize, 0)
	}
	defer rows.Close()

	var totalItems int
	items := make([]modelartwork.Summary, 0)
	for rows.Next() {
		var (
			s              modelartwork.Summary
			ownerAvatar    sql.NullString
			publishedAt    sql.NullTime
			ownerID        string
			ownerName      string
			previewImage   sql.NullString
		)
		if err := rows.Scan(
			&s.ID, &s.Title, &s.Status, &previewImage,
			&s.AvailabilityForExchange, &s.CreatedAt, &publishedAt,
			&ownerID, &ownerName, &ownerAvatar,
			&s.LikeCount, &totalItems,
		); err != nil {
			return nil, paginationMeta(page, pageSize, 0)
		}
		s.Owner = modelartwork.Owner{ID: ownerID, DisplayName: ownerName, AvatarURL: nullString(ownerAvatar)}
		s.PreviewImageURL = nullString(previewImage)
		s.PublishedAt = nullTime(publishedAt)
		items = append(items, s)
	}
	if err := rows.Err(); err != nil {
		return nil, paginationMeta(page, pageSize, 0)
	}
	return items, paginationMeta(page, pageSize, totalItems)
}

func (r *Repository) Create(input modelartwork.CreateInput) (modelartwork.Detail, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return modelartwork.Detail{}, err
	}
	defer tx.Rollback()

	const insertArtwork = `
		INSERT INTO artworks (artist_id, title, description, status, tags, preview_image_url, availability_for_exchange)
		VALUES ($1, $2, $3, 'draft', $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	var (
		artworkID         string
		createdAt         time.Time
		updatedAt         time.Time
		previewImageURL   = previewImageURL(input.Images)
	)
	if err := tx.QueryRow(
		insertArtwork,
		input.ArtistID,
		input.Title,
		input.Description,
		pq.Array(input.Tags),
		nullableString(previewImageURL),
		input.AvailabilityForExchange,
	).Scan(&artworkID, &createdAt, &updatedAt); err != nil {
		return modelartwork.Detail{}, err
	}

	if err := insertImages(tx, artworkID, input.Images); err != nil {
		return modelartwork.Detail{}, err
	}

	if err := tx.Commit(); err != nil {
		return modelartwork.Detail{}, err
	}

	return r.FindByID(artworkID)
}

func (r *Repository) FindByID(id string) (modelartwork.Detail, error) {
	const selectArtwork = `
		SELECT
			a.id, a.title, a.description, a.status, a.tags,
			a.preview_image_url, a.availability_for_exchange,
			a.created_at, a.updated_at, a.published_at,
			ar.id, ar.display_name, ar.avatar_url,
			(SELECT COUNT(*) FROM likes l WHERE l.artwork_id = a.id) AS like_count
		FROM artworks a
		JOIN artists ar ON ar.id = a.artist_id
		WHERE a.id = $1
	`

	var (
		detail        modelartwork.Detail
		previewImage  sql.NullString
		publishedAt   sql.NullTime
		ownerID       string
		ownerName     string
		ownerAvatar   sql.NullString
		tags          pq.StringArray
	)
	err := r.db.QueryRow(selectArtwork, id).Scan(
		&detail.ID, &detail.Title, &detail.Description, &detail.Status, &tags,
		&previewImage, &detail.AvailabilityForExchange,
		&detail.CreatedAt, &detail.UpdatedAt, &publishedAt,
		&ownerID, &ownerName, &ownerAvatar,
		&detail.LikeCount,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return modelartwork.Detail{}, repoerror.ErrNotFound
	}
	if err != nil {
		return modelartwork.Detail{}, err
	}
	detail.Tags = []string(tags)
	detail.Owner = modelartwork.Owner{ID: ownerID, DisplayName: ownerName, AvatarURL: nullString(ownerAvatar)}
	detail.PreviewImageURL = nullString(previewImage)
	detail.PublishedAt = nullTime(publishedAt)

	images, err := selectImages(r.db, id)
	if err != nil {
		return modelartwork.Detail{}, err
	}
	detail.Images = images

	return detail, nil
}

func (r *Repository) Update(input modelartwork.UpdateInput) (modelartwork.Detail, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return modelartwork.Detail{}, err
	}
	defer tx.Rollback()

	const updateQuery = `
		UPDATE artworks
		SET
			title = COALESCE($2, title),
			description = COALESCE($3, description),
			tags = COALESCE($4, tags),
			availability_for_exchange = COALESCE($5, availability_for_exchange),
			status = COALESCE($6, status),
			published_at = CASE
				WHEN $6 = 'published' AND published_at IS NULL THEN NOW()
				ELSE published_at
			END,
			updated_at = NOW()
		WHERE id = $1
		RETURNING id
	`

	var tags interface{}
	if input.Tags != nil {
		tags = pq.Array(input.Tags)
	}
	var status interface{}
	if input.Status != nil {
		status = string(*input.Status)
	}

	var id string
	err = tx.QueryRow(
		updateQuery,
		input.ID,
		input.Title,
		input.Description,
		tags,
		input.AvailabilityForExchange,
		status,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return modelartwork.Detail{}, repoerror.ErrNotFound
	}
	if err != nil {
		return modelartwork.Detail{}, err
	}

	if input.Images != nil {
		if _, err := tx.Exec(`DELETE FROM artwork_images WHERE artwork_id = $1`, id); err != nil {
			return modelartwork.Detail{}, err
		}
		if err := insertImages(tx, id, input.Images); err != nil {
			return modelartwork.Detail{}, err
		}
		preview := previewImageURL(input.Images)
		if _, err := tx.Exec(`UPDATE artworks SET preview_image_url = $2 WHERE id = $1`, id, nullableString(preview)); err != nil {
			return modelartwork.Detail{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return modelartwork.Detail{}, err
	}
	return r.FindByID(id)
}

func (r *Repository) Delete(id string) error {
	const query = `DELETE FROM artworks WHERE id = $1`

	result, err := r.db.Exec(query, id)
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

func insertImages(tx *sql.Tx, artworkID string, images []modelartwork.Image) error {
	const stmt = `INSERT INTO artwork_images (artwork_id, url, alt_text, position) VALUES ($1, $2, $3, $4)`
	for index, image := range images {
		if _, err := tx.Exec(stmt, artworkID, image.URL, image.AltText, index); err != nil {
			return err
		}
	}
	return nil
}

func selectImages(db *sql.DB, artworkID string) ([]modelartwork.Image, error) {
	const stmt = `SELECT url, alt_text FROM artwork_images WHERE artwork_id = $1 ORDER BY position`
	rows, err := db.Query(stmt, artworkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := make([]modelartwork.Image, 0)
	for rows.Next() {
		var (
			image   modelartwork.Image
			altText sql.NullString
		)
		if err := rows.Scan(&image.URL, &altText); err != nil {
			return nil, err
		}
		image.AltText = nullString(altText)
		images = append(images, image)
	}
	return images, rows.Err()
}

func previewImageURL(images []modelartwork.Image) *string {
	if len(images) == 0 {
		return nil
	}
	return &images[0].URL
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func nullString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	v := value.String
	return &v
}

func nullTime(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	v := value.Time
	return &v
}

func nullableString(value *string) interface{} {
	if value == nil {
		return nil
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
