package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Example_shortenURL() {
	originalURL := "https://practicum.yandex.ru"
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader(originalURL))
	if err != nil {
		fmt.Printf("Ошибка при создании запроса: %v\n", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка при чтении ответа: %v\n", err)
		return
	}

	fmt.Printf("Статус: %d\n", resp.StatusCode)
	fmt.Printf("Сокращенный URL: %s\n", string(body))

	// Output:
	// Статус: 201
	// Сокращенный URL: http://localhost:8080/abcdefgh
}

func Example_shortenURLJSON() {
	requestBody := map[string]string{
		"url": "https://practicum.yandex.ru",
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Ошибка при маршалинге JSON: %v\n", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Ошибка при создании запроса: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var response struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("Ошибка при декодировании ответа: %v\n", err)
		return
	}

	fmt.Printf("Статус: %d\n", resp.StatusCode)
	fmt.Printf("Сокращенный URL: %s\n", response.Result)

	// Output:
	// Статус: 201
	// Сокращенный URL: http://localhost:8080/abcdefgh
}

func Example_getUserURLs() {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/user/urls", nil)
	if err != nil {
		fmt.Printf("Ошибка при создании запроса: %v\n", err)
		return
	}

	req.AddCookie(&http.Cookie{
		Name:  "user_id",
		Value: "test-user-id",
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var response []struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("Ошибка при декодировании ответа: %v\n", err)
		return
	}

	fmt.Printf("Статус: %d\n", resp.StatusCode)
	fmt.Printf("Количество URL: %d\n", len(response))
	if len(response) > 0 {
		fmt.Printf("Пример URL: %s -> %s\n", response[0].ShortURL, response[0].OriginalURL)
	}

	// Output:
	// Статус: 200
	// Количество URL: 1
	// Пример URL: http://localhost:8080/abcdefgh -> https://practicum.yandex.ru
}

func Example_shortenBatch() {
	requestBody := []map[string]string{
		{
			"correlation_id": "1",
			"original_url":   "https://practicum.yandex.ru",
		},
		{
			"correlation_id": "2",
			"original_url":   "https://ya.ru",
		},
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Ошибка при маршалинге JSON: %v\n", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten/batch", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Ошибка при создании запроса: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	req.AddCookie(&http.Cookie{
		Name:  "user_id",
		Value: "test-user-id",
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var response []struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("Ошибка при декодировании ответа: %v\n", err)
		return
	}

	fmt.Printf("Статус: %d\n", resp.StatusCode)
	fmt.Printf("Количество сокращенных URL: %d\n", len(response))
	for _, item := range response {
		fmt.Printf("ID: %s, URL: %s\n", item.CorrelationID, item.ShortURL)
	}

	// Output:
	// Статус: 201
	// Количество сокращенных URL: 2
	// ID: 1, URL: http://localhost:8080/abcdefgh
	// ID: 2, URL: http://localhost:8080/ijklmnop
}

func Example_deleteURLs() {
	requestBody := []string{"abcdefgh", "ijklmnop"}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Ошибка при маршалинге JSON: %v\n", err)
		return
	}

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/user/urls", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Ошибка при создании запроса: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	req.AddCookie(&http.Cookie{
		Name:  "user_id",
		Value: "test-user-id",
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Статус: %d\n", resp.StatusCode)

	// Output:
	// Статус: 202
}
