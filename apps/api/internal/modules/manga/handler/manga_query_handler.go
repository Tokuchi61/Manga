package handler

import (
	"net/http"
	"strconv"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/manga/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) ListManga(w http.ResponseWriter, r *http.Request) {
	req := dto.ListMangaRequest{
		Search:         r.URL.Query().Get("search"),
		Genre:          r.URL.Query().Get("genre"),
		Tag:            r.URL.Query().Get("tag"),
		Theme:          r.URL.Query().Get("theme"),
		ContentWarning: r.URL.Query().Get("content_warning"),
		SortBy:         r.URL.Query().Get("sort"),
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if parsed, err := strconv.Atoi(limit); err == nil {
			req.Limit = parsed
		}
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		if parsed, err := strconv.Atoi(offset); err == nil {
			req.Offset = parsed
		}
	}

	res, err := h.service.ListPublicManga(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) GetMangaDetail(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetPublicMangaDetail(r.Context(), dto.GetMangaDetailRequest{
		MangaID: chi.URLParam(r, "manga_id"),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) ListDiscovery(w http.ResponseWriter, r *http.Request) {
	req := dto.DiscoveryRequest{
		Mode:          r.URL.Query().Get("mode"),
		CollectionKey: r.URL.Query().Get("collection"),
		SortBy:        r.URL.Query().Get("sort"),
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if parsed, err := strconv.Atoi(limit); err == nil {
			req.Limit = parsed
		}
	}

	res, err := h.service.ListDiscovery(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}
