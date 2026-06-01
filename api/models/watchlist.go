package models

import "time"

// WatchlistStatus represents the status of an anime in a watchlist.
type WatchlistStatus string

const (
	WatchlistStatusPlanToWatch WatchlistStatus = "plan_to_watch"
	WatchlistStatusWatching    WatchlistStatus = "watching"
	WatchlistStatusCompleted   WatchlistStatus = "completed"
	WatchlistStatusDropped     WatchlistStatus = "dropped"
	WatchlistStatusOnHold      WatchlistStatus = "on_hold"
)

// ValidWatchlistStatuses is the set of allowed status values.
var ValidWatchlistStatuses = map[WatchlistStatus]string{
	WatchlistStatusPlanToWatch: "Plan to Watch",
	WatchlistStatusWatching:    "Watching",
	WatchlistStatusCompleted:   "Completed",
	WatchlistStatusDropped:     "Dropped",
	WatchlistStatusOnHold:      "On Hold",
}

// WatchlistEntry represents a user's watchlist entry for an anime.
type WatchlistEntry struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	UserID    uint            `json:"user_id" gorm:"not null;index:idx_watchlist_user_anime,unique"`
	AnimeID   uint            `json:"anime_id" gorm:"not null;index:idx_watchlist_user_anime,unique;index:idx_watchlist_anime"`
	Status    WatchlistStatus `json:"status" gorm:"not null;type:varchar(20);default:'plan_to_watch'"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Anime     Anime           `json:"anime,omitempty" gorm:"foreignKey:AnimeID"`
	User      User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (WatchlistEntry) TableName() string { return "watchlists" }

// WatchlistUpsertRequest is the payload for adding/updating a watchlist entry.
type WatchlistUpsertRequest struct {
	Status WatchlistStatus `json:"status" validate:"required"`
}

// WatchlistEntryResponse is the safe response shape (no user nesting).
type WatchlistEntryResponse struct {
	ID        uint            `json:"id"`
	UserID    uint            `json:"user_id"`
	AnimeID   uint            `json:"anime_id"`
	Status    WatchlistStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Anime     *AnimeResponse  `json:"anime,omitempty"`
}

func (e WatchlistEntry) ToResponse() WatchlistEntryResponse {
	resp := WatchlistEntryResponse{
		ID: e.ID, UserID: e.UserID, AnimeID: e.AnimeID,
		Status: e.Status, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt,
	}
	if e.Anime.ID != 0 {
		ar := e.Anime.ToResponse()
		resp.Anime = &ar
	}
	return resp
}
