// Package jikan provides an HTTP client for the Jikan v4 API (https://jikan.moe/),
// an unofficial REST wrapper for MyAnimeList.  No authentication is required.
// Rate limit: ~3 requests/second — callers are responsible for adding delays.
package jikan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const BaseURL = "https://api.jikan.moe/v4"

// Client is a minimal Jikan v4 HTTP client.
type Client struct {
	http    *http.Client
	baseURL string
}

// NewClient creates a Client with a 15-second timeout.
func NewClient() *Client {
	return &Client{
		http:    &http.Client{Timeout: 15 * time.Second},
		baseURL: BaseURL,
	}
}

// -----------------------------------------------------------------
// Response types
// -----------------------------------------------------------------

// TopAnimeResponse is the envelope returned by GET /top/anime.
type TopAnimeResponse struct {
	Data       []AnimeEntry `json:"data"`
	Pagination Pagination   `json:"pagination"`
}

// Pagination describes the paging metadata from Jikan.
type Pagination struct {
	LastVisiblePage int  `json:"last_visible_page"`
	HasNextPage     bool `json:"has_next_page"`
}

// AnimeEntry is a single anime record from Jikan.
type AnimeEntry struct {
	MALID    int      `json:"mal_id"`
	Title    string   `json:"title"`
	Synopsis string   `json:"synopsis"`
	Episodes int      `json:"episodes"`
	Status   string   `json:"status"`
	Score    float64  `json:"score"`
	Aired    Aired    `json:"aired"`
	Genres   []Tag    `json:"genres"`
	Themes   []Tag    `json:"themes"`
	Images   Images   `json:"images"`
}

// Tag is a generic name-bearing object (genres, themes, demographics, …).
type Tag struct {
	Name string `json:"name"`
}

// Aired holds nullable start/end air dates in RFC3339 format.
type Aired struct {
	From *string `json:"from"`
	To   *string `json:"to"`
}

// Images contains image sets in different formats.
type Images struct {
	JPG ImageSet `json:"jpg"`
}

// ImageSet holds small and large image URLs.
type ImageSet struct {
	ImageURL      string `json:"image_url"`
	LargeImageURL string `json:"large_image_url"`
}

// -----------------------------------------------------------------
// API methods
// -----------------------------------------------------------------

// GetTopAnime fetches one page (25 entries) of the Jikan top-anime list.
// Page is 1-indexed.
func (c *Client) GetTopAnime(ctx context.Context, page int) (*TopAnimeResponse, error) {
	url := fmt.Sprintf("%s/top/anime?limit=25&page=%d", c.baseURL, page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// OK — fall through
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("jikan rate-limited (429)")
	default:
		return nil, fmt.Errorf("jikan returned HTTP %d", resp.StatusCode)
	}

	var result TopAnimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}
