package models

// ShortenRequest представляет запрос на сокращение URL
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse представляет ответ с сокращенным URL
type ShortenResponse struct {
	Result string `json:"result"`
}

// BatchRequestItem представляет элемент запроса для пакетного сокращения URL
type BatchRequestItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchResponseItem представляет элемент ответа для пакетного сокращения URL
type BatchResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
