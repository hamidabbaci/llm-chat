package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// chunkMarkdown سند رو بر اساس هدرهای ## به تکه‌های جدا تقسیم می‌کنه
func chunkMarkdown(content string) []string {
	// با این الگو، هر جا یک هدر ## شروع بشه، یک تکه‌ی جدید می‌سازیم
	re := regexp.MustCompile(`(?m)^## `)
	indexes := re.FindAllStringIndex(content, -1)

	var chunks []string
	for i, loc := range indexes {
		start := loc[0]
		end := len(content)
		if i+1 < len(indexes) {
			end = indexes[i+1][0]
		}
		chunk := strings.TrimSpace(content[start:end])
		if chunk != "" {
			chunks = append(chunks, chunk)
		}
	}
	return chunks
}

func loadKnowledgeFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return chunkMarkdown(string(data)), nil
}
func ingestFile(path string) (int, error) {
	chunks, err := loadKnowledgeFile(path)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, chunk := range chunks {
		embedding, err := getEmbedding(chunk)
		if err != nil {
			return count, err
		}

		// تبدیل []float64 به فرمتی که pgvector قبول می‌کنه، مثلا: [0.1,0.2,0.3]
		embeddingStr := floatsToVectorString(embedding)

		_, err = db.Exec(
			"INSERT INTO knowledge (content, embedding) VALUES ($1, $2)",
			chunk, embeddingStr,
		)
		if err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func floatsToVectorString(values []float64) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = fmt.Sprintf("%f", v)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

type searchResult struct {
	Content  string  `json:"content"`
	Distance float64 `json:"distance"`
}

func searchKnowledge(query string, limit int) ([]searchResult, error) {
	queryEmbedding, err := getEmbedding(query)
	if err != nil {
		return nil, err
	}
	queryVector := floatsToVectorString(queryEmbedding)

	rows, err := db.Query(
		`SELECT content, embedding <=> $1 AS distance
		 FROM knowledge
		 ORDER BY distance ASC
		 LIMIT $2`,
		queryVector, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []searchResult
	for rows.Next() {
		var r searchResult
		if err := rows.Scan(&r.Content, &r.Distance); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}
