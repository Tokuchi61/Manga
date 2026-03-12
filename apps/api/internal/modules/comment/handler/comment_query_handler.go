package handler

import (
	"net/http"
	"strconv"

	"github.com/Tokuchi61/Manga/apps/api/internal/modules/comment/dto"
	"github.com/go-chi/chi/v5"
)

func (h *HTTPHandler) ListByTarget(w http.ResponseWriter, r *http.Request) {
	req := dto.ListCommentsRequest{
		TargetType: chi.URLParam(r, "target_type"),
		TargetID:   chi.URLParam(r, "target_id"),
		SortBy:     r.URL.Query().Get("sort"),
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
	if includeHidden := r.URL.Query().Get("include_hidden"); includeHidden != "" {
		if parsed, err := strconv.ParseBool(includeHidden); err == nil {
			req.IncludeHidden = parsed
		}
	}

	res, err := h.service.ListCommentsByTarget(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) GetDetail(w http.ResponseWriter, r *http.Request) {
	req := dto.GetCommentDetailRequest{CommentID: chi.URLParam(r, "comment_id")}
	if includeHidden := r.URL.Query().Get("include_hidden"); includeHidden != "" {
		if parsed, err := strconv.ParseBool(includeHidden); err == nil {
			req.IncludeHidden = parsed
		}
	}

	res, err := h.service.GetCommentDetail(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *HTTPHandler) GetThread(w http.ResponseWriter, r *http.Request) {
	req := dto.GetCommentThreadRequest{
		CommentID: chi.URLParam(r, "comment_id"),
		SortBy:    r.URL.Query().Get("sort"),
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
	if includeHidden := r.URL.Query().Get("include_hidden"); includeHidden != "" {
		if parsed, err := strconv.ParseBool(includeHidden); err == nil {
			req.IncludeHidden = parsed
		}
	}

	res, err := h.service.GetCommentThread(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, res)
}
