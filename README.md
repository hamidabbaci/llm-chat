# 🤖 چت‌بات هوش مصنوعی با Gemini (Go)

یک چت‌بات ساده و امن، نوشته‌شده با **Go**، که از **Gemini API** گوگل برای تولید پاسخ استفاده می‌کنه و تاریخچه‌ی مکالمه رو در **PostgreSQL** ذخیره می‌کنه تا حتی بعد از ری‌استارت سرور هم از بین نره.

این پروژه به‌عنوان یک پروژه‌ی یادگیری ساخته شده، با تمرکز روی یادگیری اصول LLM integration، مدیریت امن کلیدها، و کار با دیتابیس در Go.

## ✨ امکانات


- رابط وب ساده و واکنش‌گرا برای چت
- اتصال به **Gemini API** برای تولید پاسخ
- نگهداری تاریخچه‌ی کامل مکالمه در **PostgreSQL** (context window واقعی، نه فقط حافظه‌ی موقت)
- زیرساخت دیتابیس با **Docker Compose** (بدون نیاز به نصب دستی Postgres)
- مدیریت امن secretها (کلید API و رمز دیتابیس هیچ‌وقت وارد Git نمی‌شن)

## 🛠 تکنولوژی‌ها

- **Go** (net/http، بدون framework اضافه)
- **PostgreSQL** (از طریق `lib/pq`)
- **Docker Compose**
- **Gemini API** (Google AI Studio)

## 🚀 راه‌اندازی

### ۱. دریافت پروژه

```bash
git clone https://github.com/hamidabbaci/llm-chat.git
cd llm-chat
```

### ۲. تنظیم متغیرهای محیطی

یک فایل `.env` بساز (از روی نمونه‌ی `.env.example`) و رمز دیتابیس دلخواهت رو توش بذار:

```bash
cp .env.example .env
```

### ۳. ساخت فایل کلیدها

یک فایل `config.go` کنار بقیه‌ی فایل‌ها بساز:

```go
package main

const apiKey = "کلید Gemini API خودت"
const dbPassword = "همون رمزی که توی .env گذاشتی"
```

> کلید رایگان Gemini API رو می‌تونی از [aistudio.google.com](https://aistudio.google.com/apikey) بگیری.

### ۴. بالا آوردن دیتابیس

```bash
docker compose up -d
```

### ۵. ساخت جدول (فقط بار اول)

```bash
docker exec -it llm-chat-db psql -U postgres -d llmchat -c "
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);"
```

### ۶. اجرای برنامه

```bash
go run main.go
```

برو به `http://localhost:8080` و شروع کن به چت کردن 🎉

## 🔒 امنیت

کلید API و رمز دیتابیس هیچ‌وقت داخل کد commit نمی‌شن:
- `config.go` → حاوی کلید واقعی API، در `.gitignore` قرار داره
- `.env` → حاوی رمز واقعی دیتابیس، در `.gitignore` قرار داره
- `.env.example` و این README به‌جاش نشون می‌دن که چه مقادیری لازمه

## 📌 نقشه‌راه آینده

- [ ] اضافه‌کردن Embeddings و Semantic Search
- [ ] اتصال به Vector Database برای RAG
- [ ] پیاده‌سازی یک AI Agent ساده

