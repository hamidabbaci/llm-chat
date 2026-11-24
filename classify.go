package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type classification struct {
	Category string `json:"category"`
	Urgency  string `json:"urgency"`
	Summary  string `json:"summary"`
}

func classifyMessage(message string) (*classification, error) {
	prompt := `تو یک دستیار دسته‌بندی پیام‌های پشتیبانی هستی.
پیام کاربر زیر رو تحلیل کن و فقط و فقط یک JSON با این ساختار دقیق برگردون، بدون هیچ توضیح اضافه، بدون ` + "```json" + `، فقط خود JSON خام:

{
  "category": "یکی از این مقادیر: deposit, withdrawal, kyc, technical, other",
  "urgency": "یکی از این مقادیر: low, medium, high",
  "summary": "خلاصه‌ی یک‌خطی مشکل به فارسی"
}

پیام کاربر: "` + message + `"`

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-3.5-flash-lite:generateContent?key=" + apiKey

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
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

	var parsed geminiResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}

	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("جوابی دریافت نشد: %s", string(body))
	}

	rawText := parsed.Candidates[0].Content.Parts[0].Text

	var result classification
	if err := json.Unmarshal([]byte(rawText), &result); err != nil {
		return nil, fmt.Errorf("جواب مدل JSON معتبر نبود: %s", rawText)
	}

	return &result, nil
}
