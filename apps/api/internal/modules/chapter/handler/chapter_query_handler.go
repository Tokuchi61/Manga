package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/chapter/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) ListByManga(w http.ResponseWriter, r *http.Request) {
	req := dto.ListChapterRequest{
		MangaID: chi.URLParam(r, "manga_id"),
		SortBy:  r.URL.Query().Get("sort"),
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

	res, err := h.service.ListChaptersByManga(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) GetDetail(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetChapterDetail(r.Context(), dto.GetChapterDetailRequest{ChapterID: chi.URLParam(r, "chapter_id")})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) Read(w http.ResponseWriter, r *http.Request) {
	req := dto.ReadChapterRequest{
		ChapterID: chi.URLParam(r, "chapter_id"),
		Mode:      r.URL.Query().Get("mode"),
	}
	if at := r.URL.Query().Get("at"); at != "" {
		parsed, err := time.Parse(time.RFC3339, at)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		req.At = &parsed
	}

	res, err := h.service.ReadChapter(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) Navigation(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetNavigation(r.Context(), dto.NavigationRequest{ChapterID: chi.URLParam(r, "chapter_id")})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}
