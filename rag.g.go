package main

import (
	"fmt"
	"strings"
)

func answerWithRAG(question string) (string, error) {
	results, err := searchKnowledge(question, 3)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "اطلاعاتی برای پاسخ به این سؤال پیدا نشد.", nil
	}

	var contextParts []string
	for _, r := range results {
		contextParts = append(contextParts, r.Content)
	}
	context := strings.Join(contextParts, "\n\n---\n\n")

	prompt := fmt.Sprintf(`تو پشتیبان هوش مصنوعی صرافی Fund Wayne هستی.
فقط و فقط از روی اطلاعات زیر به سؤال کاربر جواب بده. اگه جواب توی این اطلاعات نبود، بگو «این سؤال نیاز به بررسی توسط سرپرست داره» و چیزی از خودت اختراع نکن.

اطلاعات:
%s

سؤال کاربر: %s`, context, question)

	return askGeminiSimple(prompt)
}

func askGeminiSimple(prompt string) (string, error) {
	history := []map[string]interface{}{
		{
			"role":  "user",
			"parts": []map[string]string{{"text": prompt}},
		},
	}
	return askGemini(history)
}
