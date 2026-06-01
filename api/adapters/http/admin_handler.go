package httphandler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"myanimeapi/api/middleware"
	"myanimeapi/api/models"
	"myanimeapi/internal/db"
	"myanimeapi/internal/errors"

	"github.com/gorilla/mux"
)

// AdminHandler handles admin-only HTTP endpoints.
type AdminHandler struct {
	db db.DBInterface
}

// NewAdminHandler creates a new AdminHandler with the given database interface.
func NewAdminHandler(database db.DBInterface) *AdminHandler {
	return &AdminHandler{db: database}
}

// AdminStats represents aggregate counts for the dashboard.
type AdminStats struct {
	UsersTotal     int64 `json:"users_total"`
	AnimeTotal     int64 `json:"anime_total"`
	WithCoverTotal int64 `json:"with_cover_total"`
	ReviewsTotal   int64 `json:"reviews_total"`
	FavoritesTotal int64 `json:"favorites_total"`
	GenresTotal    int64 `json:"genres_total"`
	TagsTotal      int64 `json:"tags_total"`
}

// StatsHandler godoc
// @Summary Admin dashboard stats
// @Tags admin
// @Produce json
// @Success 200 {object} AdminStats
// @Security BearerAuth
// @Router /admin/stats [get]
func (h *AdminHandler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	gdb := h.db.WithContext(r.Context())
	var stats AdminStats
	gdb.Model(&models.User{}).Count(&stats.UsersTotal)
	gdb.Model(&models.Anime{}).Count(&stats.AnimeTotal)
	gdb.Model(&models.Anime{}).Where("cover_url IS NOT NULL AND cover_url != ''").Count(&stats.WithCoverTotal)
	gdb.Model(&models.Review{}).Count(&stats.ReviewsTotal)
	gdb.Model(&models.Favorite{}).Count(&stats.FavoritesTotal)
	gdb.Model(&models.Genre{}).Count(&stats.GenresTotal)
	gdb.Model(&models.Tag{}).Count(&stats.TagsTotal)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// AdminUserRow is the safe subset of User returned to the admin panel.
type AdminUserRow struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsAdmin   bool      `json:"is_admin"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// UsersHandler lists all users with pagination.
// @Summary List all users (admin)
// @Tags admin
// @Produce json
// @Param page query int false "Page (default 1)"
// @Param limit query int false "Page size (default 20, max 100)"
// @Security BearerAuth
// @Router /admin/users [get]
func (h *AdminHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	page, limit := 1, 20
	if v := r.URL.Query().Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	gdb := h.db.WithContext(r.Context())
	var total int64
	gdb.Model(&models.User{}).Count(&total)

	var users []models.User
	gdb.Select("id, username, email, role, is_admin, is_active, created_at").
		Order("created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&users)

	rows := make([]AdminUserRow, len(users))
	for i, u := range users {
		rows[i] = AdminUserRow{
			ID: u.ID, Username: u.Username, Email: u.Email,
			Role: u.Role, IsAdmin: u.IsAdmin, IsActive: u.IsActive,
			CreatedAt: u.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": rows, "total": total, "page": page, "limit": limit,
	})
}

// UpdateRoleRequest is the payload for changing a user's role.
type UpdateRoleRequest struct {
	Role string `json:"role" validate:"required"`
}

// UpdateRoleHandler changes the role of a user.
// @Summary Update user role (admin)
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body UpdateRoleRequest true "New role"
// @Security BearerAuth
// @Router /admin/users/{id}/role [put]
func (h *AdminHandler) UpdateRoleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid user ID", "", nil)
		return
	}

	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid JSON", err.Error(), nil)
		return
	}
	validRoles := map[string]bool{"user": true, "reviewer": true, "admin": true}
	if !validRoles[req.Role] {
		errors.WriteErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidInput, "Invalid role",
			"role must be one of: user, reviewer, admin", nil)
		return
	}

	gdb := h.db.WithContext(r.Context())
	updates := map[string]interface{}{"role": req.Role, "is_admin": req.Role == "admin"}
	if res := gdb.Model(&models.User{}).Where("id = ?", userID).Updates(updates); res.Error != nil {
		errors.WriteErrorResponse(w, http.StatusInternalServerError, errors.ErrInternalServer, "Update failed", res.Error.Error(), nil)
		return
	} else if res.RowsAffected == 0 {
		errors.WriteErrorResponse(w, http.StatusNotFound, errors.ErrResourceNotFound, "User not found", "", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"role": req.Role})
}

// TopFavoritedRow holds a single row from the top-favorited analytics query.
type TopFavoritedRow struct {
	AnimeID uint   `json:"anime_id"`
	Title   string `json:"title"`
	Count   int64  `json:"count"`
}

// TopReviewedRow holds a single row from the top-reviewed analytics query.
type TopReviewedRow struct {
	AnimeID   uint    `json:"anime_id"`
	Title     string  `json:"title"`
	Count     int64   `json:"count"`
	AvgRating float64 `json:"avg_rating"`
}

// DailyCountRow holds a date and a count for time-series analytics.
type DailyCountRow struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// AnalyticsHandler returns aggregated analytics data for the admin panel.
// @Summary Admin analytics
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Router /admin/analytics [get]
func (h *AdminHandler) AnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	gdb := h.db.WithContext(r.Context())

	var topFavorited []TopFavoritedRow
	gdb.Raw(`SELECT favorites.anime_id, animes.title, COUNT(*) as count
		FROM favorites
		JOIN animes ON animes.id = favorites.anime_id
		GROUP BY favorites.anime_id, animes.title
		ORDER BY count DESC
		LIMIT 5`).Scan(&topFavorited)

	var topReviewed []TopReviewedRow
	gdb.Raw(`SELECT reviews.anime_id, animes.title, COUNT(*) as count, ROUND(AVG(reviews.rating)::numeric, 1) as avg_rating
		FROM reviews
		JOIN animes ON animes.id = reviews.anime_id
		GROUP BY reviews.anime_id, animes.title
		ORDER BY count DESC
		LIMIT 5`).Scan(&topReviewed)

	var recentReviewsPerDay []DailyCountRow
	gdb.Raw(`SELECT DATE(created_at) as date, COUNT(*) as count
		FROM reviews
		WHERE created_at >= NOW() - INTERVAL '7 days'
		GROUP BY DATE(created_at)
		ORDER BY date DESC`).Scan(&recentReviewsPerDay)

	var recentSignupsPerDay []DailyCountRow
	gdb.Raw(`SELECT DATE(created_at) as date, COUNT(*) as count
		FROM users
		WHERE created_at >= NOW() - INTERVAL '7 days'
		GROUP BY DATE(created_at)
		ORDER BY date DESC`).Scan(&recentSignupsPerDay)

	// Ensure nil slices are marshalled as empty arrays
	if topFavorited == nil {
		topFavorited = []TopFavoritedRow{}
	}
	if topReviewed == nil {
		topReviewed = []TopReviewedRow{}
	}
	if recentReviewsPerDay == nil {
		recentReviewsPerDay = []DailyCountRow{}
	}
	if recentSignupsPerDay == nil {
		recentSignupsPerDay = []DailyCountRow{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"top_favorited":          topFavorited,
		"top_reviewed":           topReviewed,
		"recent_reviews_per_day": recentReviewsPerDay,
		"recent_signups_per_day": recentSignupsPerDay,
	})
}

// RegisterAdminRoutes wires admin endpoints (all admin-only).
func (h *AdminHandler) RegisterAdminRoutes(router *mux.Router) {
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware)

	adminRouter.Handle("/stats", middleware.RequireAdmin(http.HandlerFunc(h.StatsHandler))).Methods("GET")
	adminRouter.Handle("/users", middleware.RequireAdmin(http.HandlerFunc(h.UsersHandler))).Methods("GET")
	adminRouter.Handle("/users/{id}/role", middleware.RequireAdmin(http.HandlerFunc(h.UpdateRoleHandler))).Methods("PUT")
	adminRouter.Handle("/analytics", middleware.RequireAdmin(http.HandlerFunc(h.AnalyticsHandler))).Methods("GET")
}
