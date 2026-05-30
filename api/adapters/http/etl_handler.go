package httphandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"myanimeapi/api/middleware"
	"myanimeapi/api/services"
	"myanimeapi/internal/errors"

	"github.com/gorilla/mux"
)

// ETLHandler exposes admin-only endpoints for data ingestion.
type ETLHandler struct {
	etlService services.ETLServiceInterface
}

// NewETLHandler creates an ETLHandler.
func NewETLHandler(etlService services.ETLServiceInterface) *ETLHandler {
	return &ETLHandler{etlService: etlService}
}

// SyncJikanHandler godoc
// @Summary      Sync anime data from Jikan / MyAnimeList
// @Description  Fetches the top-anime list from Jikan API (jikan.moe) and upserts
// @Description  each entry into the local database, including genres, tags, cover
// @Description  images, and MyAnimeList scores.  Admin-only.
// @Tags         etl
// @Produce      json
// @Param        pages  query  int  false  "Pages to fetch from Jikan (1 page = 25 anime, default 5, max 20)"
// @Success      200  {object}  services.SyncResult
// @Failure      401  {object}  errors.ErrorResponse
// @Failure      403  {object}  errors.ErrorResponse
// @Failure      500  {object}  errors.ErrorResponse
// @Security     BearerAuth
// @Router       /etl/sync [post]
func (h *ETLHandler) SyncJikanHandler(w http.ResponseWriter, r *http.Request) {
	pages := 5
	if q := r.URL.Query().Get("pages"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 {
			if n > 20 {
				n = 20
			}
			pages = n
		}
	}

	result, err := h.etlService.SyncFromJikan(r.Context(), pages)
	if err != nil {
		errors.WriteErrorResponse(
			w,
			http.StatusInternalServerError,
			errors.ErrInternalServer,
			"ETL sync failed",
			err.Error(),
			nil,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result) //nolint:errcheck
}

// RegisterETLRoutes wires ETL endpoints under /v1/etl (admin-only).
func (h *ETLHandler) RegisterETLRoutes(router *mux.Router) {
	etlRouter := router.PathPrefix("/etl").Subrouter()
	etlRouter.Use(middleware.AuthMiddleware)
	etlRouter.Handle(
		"/sync",
		middleware.RequireAdmin(http.HandlerFunc(h.SyncJikanHandler)),
	).Methods("POST")
}
