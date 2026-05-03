package artwork

import "time"

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
)

type Image struct {
	URL     string
	AltText *string
}

type Owner struct {
	ID          string
	DisplayName string
	AvatarURL   *string
}

type Artwork struct {
	ID                      string
	ArtistID                string
	Title                   string
	Description             string
	Status                  Status
	Tags                    []string
	Images                  []Image
	PreviewImageURL         *string
	LikeCount               int
	LikedByMe               *bool
	AvailabilityForExchange bool
	CreatedAt               time.Time
	UpdatedAt               time.Time
	PublishedAt             *time.Time
}

type Summary struct {
	ID                      string
	Title                   string
	Status                  Status
	Owner                   Owner
	PreviewImageURL         *string
	LikeCount               int
	LikedByMe               *bool
	AvailabilityForExchange bool
	CreatedAt               time.Time
	PublishedAt             *time.Time
}

type Detail struct {
	ID                      string
	Title                   string
	Description             string
	Status                  Status
	Owner                   Owner
	PreviewImageURL         *string
	LikeCount               int
	LikedByMe               *bool
	AvailabilityForExchange bool
	Tags                    []string
	Images                  []Image
	CreatedAt               time.Time
	UpdatedAt               time.Time
	PublishedAt             *time.Time
}

type OwneredRef struct {
	ID              string
	Title           string
	OwnerID         string
	PreviewImageURL *string
}

type CreateInput struct {
	ArtistID                string
	Title                   string
	Description             string
	Tags                    []string
	Images                  []Image
	AvailabilityForExchange bool
}

type UpdateInput struct {
	ID                      string
	Title                   *string
	Description             *string
	Tags                    []string
	Images                  []Image
	AvailabilityForExchange *bool
	Status                  *Status
}

type ListFilter struct {
	Query    *string
	ArtistID *string
	Status   *Status
	Sort     string
	Page     int
	PageSize int
}
