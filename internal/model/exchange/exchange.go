package exchange

import (
	"time"

	modelartwork "github.com/andrqxa/artshare/internal/model/artwork"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusAccepted  Status = "accepted"
	StatusRejected  Status = "rejected"
	StatusCancelled Status = "cancelled"
)

type Scope string

const (
	ScopeInbox  Scope = "inbox"
	ScopeOutbox Scope = "outbox"
	ScopeAll    Scope = "all"
)

type Request struct {
	ID          string
	ArtworkID   string
	RequesterID string
	OwnerID     string
	Status      Status
	Message     *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Summary struct {
	ID          string
	Artwork     modelartwork.OwneredRef
	RequesterID string
	OwnerID     string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Detail struct {
	ID          string
	Artwork     modelartwork.OwneredRef
	RequesterID string
	OwnerID     string
	Status      Status
	Message     *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateInput struct {
	ArtworkID   string
	RequesterID string
	Message     *string
}

type UpdateStatusInput struct {
	ID      string
	Status  Status
	Message *string
}

type ListFilter struct {
	UserID   string
	Scope    Scope
	Status   *Status
	Page     int
	PageSize int
}
