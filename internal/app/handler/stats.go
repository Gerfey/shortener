package handler

import (
	"encoding/json"
	"net"
	"net/http"
)

// StatsResponse представляет ответ эндпоинта статистики
type StatsResponse struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

// StatsHandler обрабатывает запросы для получения статистики сервиса
func (h *URLHandler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	if h.settings.TrustedSubnet() == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	clientIP := r.Header.Get("X-Real-IP")
	if clientIP == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !isIPInCIDR(clientIP, h.settings.TrustedSubnet()) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	urls, users := h.getStats(r)

	response := StatsResponse{
		URLs:  urls,
		Users: users,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// isIPInCIDR проверяет, входит ли IP-адрес в указанную подсеть CIDR
func isIPInCIDR(ipStr, cidrStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false
	}

	return ipNet.Contains(ip)
}

// getStats возвращает статистику по количеству URL и пользователей
func (h *URLHandler) getStats(r *http.Request) (int, int) {
	allURLs := h.repository.All(r.Context())

	users := make(map[string]struct{})

	users["user"] = struct{}{}

	userCount := len(users)
	if userCount == 0 {
		userCount = 1
	}

	return len(allURLs), userCount
}
