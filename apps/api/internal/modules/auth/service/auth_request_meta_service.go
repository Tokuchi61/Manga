package service

import "strings"

func normalizeMeta(meta RequestMeta) RequestMeta {
	meta.RequestID = strings.TrimSpace(meta.RequestID)
	meta.CorrelationID = strings.TrimSpace(meta.CorrelationID)
	if meta.CorrelationID == "" {
		meta.CorrelationID = meta.RequestID
	}
	meta.Device = strings.TrimSpace(meta.Device)
	meta.IP = strings.TrimSpace(meta.IP)
	return meta
}

func coalesce(primary string, fallback string) string {
	if strings.TrimSpace(primary) != "" {
		return strings.TrimSpace(primary)
	}
	return strings.TrimSpace(fallback)
}
