package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
)

type embeddingResponse struct {
	Embedding struct {
		Values []float64 `json:"values"`
	} `json:"embedding"`
}

func getEmbedding(text string) ([]float64, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-embedding-001:embedContent?key=" + apiKey

	requestBody := map[string]interface{}{
		"content": map[string]interface{}{
			"parts": []map[string]string{
				{"text": text},
			},
		},
		"taskType": "SEMANTIC_SIMILARITY",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed embeddingResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}

	if len(parsed.Embedding.Values) == 0 {
		return nil, fmt.Errorf("embedding دریافت نشد: %s", string(body))
	}

	return parsed.Embedding.Values, nil
} // cosineSimilarity میزان شباهت دو بردار رو بین -1 تا 1 برمی‌گردونه
// هرچی به 1 نزدیک‌تر باشه، یعنی از نظر معنایی شبیه‌ترن
func cosineSimilarity(a, b []float64) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
