package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type chatRequest struct {
	Message string `json:"message"`
}

type chatResponse struct {
	Reply string `json:"reply"`
}

func askGemini(fullHistory []map[string]interface{}) (string, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-3.5-flash-lite:generateContent?key=" + apiKey

	requestBody := map[string]interface{}{
		"contents": fullHistory,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var parsed geminiResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", err
	}

	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("جوابی دریافت نشد: %s", string(body))
	}

	return parsed.Candidates[0].Content.Parts[0].Text, nil
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "بدنه‌ی درخواست نامعتبر است", http.StatusBadRequest)
		return
	}

	// پیام کاربر رو ذخیره کن
	if err := saveMessage("user", req.Message); err != nil {
		http.Error(w, "خطا در ذخیره‌ی پیام: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// کل تاریخچه رو از دیتابیس بخون
	history, err := loadHistory()
	if err != nil {
		http.Error(w, "خطا در خواندن تاریخچه: "+err.Error(), http.StatusInternalServerError)
		return
	}

	answer, err := askGemini(history)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// جواب مدل رو هم ذخیره کن
	if err := saveMessage("model", answer); err != nil {
		http.Error(w, "خطا در ذخیره‌ی جواب: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatResponse{Reply: answer})
}

func main() {
	initDB()

	http.HandleFunc("/chat", chatHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("سرور روی http://localhost:8080 اجراست")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("خطا در اجرای سرور:", err)
	}
}
