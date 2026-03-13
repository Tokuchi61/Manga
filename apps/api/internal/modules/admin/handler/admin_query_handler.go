package handler

import (
	"net/http"
	"strconv"
	"strings"
)

func parseIntQuery(r *http.Request, key string) int {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return 0
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return value
}

func parseBoolQuery(r *http.Request, key string) bool {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return false
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false
	}
	return value
}

func parseStringQuery(r *http.Request, key string) string {
	return strings.TrimSpace(r.URL.Query().Get(key))
}
