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
	clientIP := r.Header.Get("X-Real-IP")
	if clientIP == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	stats, err := h.statsUseCase.GetStats(r.Context(), clientIP)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	response := StatsResponse{
		URLs:  stats.URLs,
		Users: stats.Users,
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


