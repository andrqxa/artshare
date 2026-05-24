//go:build integration

package artist_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const (
	defaultBaseURL     = "http://localhost:8080"
	defaultDatabaseURL = "postgres://artshare:artshare@localhost:5432/artshare?sslmode=disable"
	apiPrefix          = "/api/v1"
)

type artistDetail struct {
	ID          string  `json:"id"`
	DisplayName string  `json:"displayName"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	UserID      string  `json:"userId"`
}

type artistSummary struct {
	ID          string  `json:"id"`
	DisplayName string  `json:"displayName"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
	Bio         *string `json:"bio,omitempty"`
}

type artistListResponse struct {
	Data []artistSummary `json:"data"`
	Meta struct {
		Page       int `json:"page"`
		PageSize   int `json:"pageSize"`
		TotalItems int `json:"totalItems"`
		TotalPages int `json:"totalPages"`
	} `json:"meta"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type fixture struct {
	t       *testing.T
	baseURL string
	client  *http.Client
	db      *sql.DB
}

func newFixture(t *testing.T) *fixture {
	t.Helper()

	baseURL := envOr("ARTSHARE_BASE_URL", defaultBaseURL)
	databaseURL := envOr("DATABASE_URL", defaultDatabaseURL)

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	fx := &fixture{
		t:       t,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
		db:      db,
	}

	fx.waitReady()
	fx.cleanup()

	return fx
}

func (f *fixture) waitReady() {
	f.t.Helper()

	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := f.client.Get(f.baseURL + apiPrefix + "/artists?pageSize=1")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 500 {
				return
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	f.t.Fatalf("service at %s did not become ready", f.baseURL)
}

func (f *fixture) cleanup() {
	f.t.Helper()

	statements := []string{
		"DELETE FROM exchange_requests",
		"DELETE FROM likes",
		"DELETE FROM artwork_images",
		"DELETE FROM artworks",
		"DELETE FROM artists",
		"DELETE FROM users",
	}
	for _, stmt := range statements {
		if _, err := f.db.ExecContext(context.Background(), stmt); err != nil {
			f.t.Fatalf("cleanup %q: %v", stmt, err)
		}
	}
}

func (f *fixture) seedUser(t *testing.T, email string) string {
	t.Helper()

	const stmt = `
		INSERT INTO users (email, password, role)
		VALUES ($1, 'password', 'artist')
		RETURNING id
	`
	var id string
	if err := f.db.QueryRow(stmt, email).Scan(&id); err != nil {
		t.Fatalf("seed user %q: %v", email, err)
	}
	return id
}

func (f *fixture) do(t *testing.T, method, path string, userID string, body any) (*http.Response, []byte) {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reqBody = bytes.NewReader(buf)
	}

	req, err := http.NewRequest(method, f.baseURL+path, reqBody)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if userID != "" {
		req.Header.Set("Authorization", "Bearer dev-token-"+userID)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		t.Fatalf("perform request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}
	return resp, bodyBytes
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func TestArtist_CreateAndGetCurrent(t *testing.T) {
	fx := newFixture(t)

	userID := fx.seedUser(t, "create@artshare.test")

	bio := "loves drawing seascapes"
	avatar := "https://example.test/a.png"
	resp, body := fx.do(t, http.MethodPost, apiPrefix+"/artists/me", userID, map[string]any{
		"displayName": "Marina",
		"bio":         bio,
		"avatarUrl":   avatar,
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create artist: status=%d body=%s", resp.StatusCode, body)
	}

	var created artistDetail
	if err := json.Unmarshal(body, &created); err != nil {
		t.Fatalf("decode created: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected non-empty id, got %q", created.ID)
	}
	if created.DisplayName != "Marina" {
		t.Fatalf("display name: got %q want Marina", created.DisplayName)
	}
	if created.UserID != userID {
		t.Fatalf("user id: got %q want %q", created.UserID, userID)
	}

	resp, body = fx.do(t, http.MethodGet, apiPrefix+"/artists/me", userID, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get current: status=%d body=%s", resp.StatusCode, body)
	}
	var fetched artistDetail
	if err := json.Unmarshal(body, &fetched); err != nil {
		t.Fatalf("decode fetched: %v", err)
	}
	if fetched.ID != created.ID {
		t.Fatalf("fetched id: got %q want %q", fetched.ID, created.ID)
	}
}

func TestArtist_CreateConflict(t *testing.T) {
	fx := newFixture(t)

	userID := fx.seedUser(t, "conflict@artshare.test")

	resp, body := fx.do(t, http.MethodPost, apiPrefix+"/artists/me", userID, map[string]any{
		"displayName": "Once",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("first create: status=%d body=%s", resp.StatusCode, body)
	}

	resp, body = fx.do(t, http.MethodPost, apiPrefix+"/artists/me", userID, map[string]any{
		"displayName": "Twice",
	})
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 on duplicate, got status=%d body=%s", resp.StatusCode, body)
	}

	var er errorResponse
	if err := json.Unmarshal(body, &er); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if er.Code == 0 {
		t.Fatalf("expected non-zero error code in %s", body)
	}
}

func TestArtist_CreateValidation(t *testing.T) {
	fx := newFixture(t)

	userID := fx.seedUser(t, "val@artshare.test")

	resp, body := fx.do(t, http.MethodPost, apiPrefix+"/artists/me", userID, map[string]any{
		"displayName": "x",
	})
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 on short displayName, got status=%d body=%s", resp.StatusCode, body)
	}
	if !strings.Contains(string(body), "displayName") {
		t.Fatalf("expected validation message about displayName, got %s", body)
	}
}

func TestArtist_Update(t *testing.T) {
	fx := newFixture(t)

	userID := fx.seedUser(t, "update@artshare.test")
	fx.do(t, http.MethodPost, apiPrefix+"/artists/me", userID, map[string]any{
		"displayName": "Original",
	})

	newName := "Updated"
	resp, body := fx.do(t, http.MethodPatch, apiPrefix+"/artists/me", userID, map[string]any{
		"displayName": newName,
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update: status=%d body=%s", resp.StatusCode, body)
	}
	var updated artistDetail
	if err := json.Unmarshal(body, &updated); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if updated.DisplayName != newName {
		t.Fatalf("display name: got %q want %q", updated.DisplayName, newName)
	}
}

func TestArtist_GetByIDNotFound(t *testing.T) {
	fx := newFixture(t)

	resp, body := fx.do(t, http.MethodGet, apiPrefix+"/artists/00000000-0000-0000-0000-000000000999", "", nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got status=%d body=%s", resp.StatusCode, body)
	}
}

func TestArtist_List(t *testing.T) {
	fx := newFixture(t)

	for i, name := range []string{"Alice", "Bob", "Carol"} {
		userID := fx.seedUser(t, fmt.Sprintf("list-%d@artshare.test", i))
		resp, body := fx.do(t, http.MethodPost, apiPrefix+"/artists/me", userID, map[string]any{
			"displayName": name,
		})
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("seed artist %q: status=%d body=%s", name, resp.StatusCode, body)
		}
	}

	resp, body := fx.do(t, http.MethodGet, apiPrefix+"/artists?pageSize=10", "", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list: status=%d body=%s", resp.StatusCode, body)
	}
	var list artistListResponse
	if err := json.Unmarshal(body, &list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if list.Meta.TotalItems != 3 {
		t.Fatalf("total items: got %d want 3", list.Meta.TotalItems)
	}
	if len(list.Data) != 3 {
		t.Fatalf("data length: got %d want 3", len(list.Data))
	}
	if list.Data[0].DisplayName != "Alice" {
		t.Fatalf("expected alphabetic ordering by displayName, got %q first", list.Data[0].DisplayName)
	}

	q := url.Values{}
	q.Set("q", "Bob")
	resp, body = fx.do(t, http.MethodGet, apiPrefix+"/artists?"+q.Encode(), "", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list with query: status=%d body=%s", resp.StatusCode, body)
	}
	var filtered artistListResponse
	if err := json.Unmarshal(body, &filtered); err != nil {
		t.Fatalf("decode filtered list: %v", err)
	}
	if filtered.Meta.TotalItems != 1 || len(filtered.Data) != 1 || filtered.Data[0].DisplayName != "Bob" {
		t.Fatalf("expected single Bob result, got %+v", filtered)
	}
}
