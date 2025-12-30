## 1. Teknologi yang Digunakan

- **Bahasa Pemrograman:** Go (Golang)
- **Database:** Supabase (PostgreSQL)
- **Router / Framework HTTP:** Gin
- **Tool / Library Tambahan:**
  - `github.com/gin-gonic/gin` → routing & HTTP framework
  - `github.com/joho/godotenv` → load environment variables
  - Supabase REST API (PostgREST)

---

## 2. Tujuan Project

Project ini dibuat untuk:

- Membuat API backend sederhana untuk manajemen **user** dan **task**.
- Menyediakan fitur CRUD (Create, Read, Update, Delete) untuk task.
- Mendukung relasi user → task melalui endpoint nested.
- Menerapkan error handling standar (**400, 404, 500**).
- Menggunakan Go dan Supabase sebagai backend API.

Project ini dikerjakan **secara individu**.

---

## 3. Struktur Project

project/
├─ database/
│ └─ db.go
├─ handlers/
│ ├─ user.go
│ └─ task.go
├─ models/
│ ├─ user.go
│ └─ task.go
├─ .env
├─ go.mod
├─ main.go


## 4. Cara Install dan Menjalankan

### Clone Repository
git clone https://github.com/BettingBalls/BackendApi-Go.git
cd BackendApi-Go

# Install dependencies
go mod tidy

# Jalankan server
go run main.go

## 5 Contoh Request dan Response API
BASE_URL="http://localhost:8080"

================== TASK ==================
# 1. GET /tasks
curl -s -X GET "$BASE_URL/tasks" | jq
echo -e "\n"

# 2. GET /users/{id}/tasks
curl -s -X GET "$BASE_URL/users/1/tasks" | jq
echo -e "\n"

# 3. POST /tasks
curl -s -X POST "$BASE_URL/tasks" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "user_id": 1,
    "title": "Belajar Go",
    "description": "Testing API",
    "status": "done",
    "deadline": "2025-12-31"
  }' | jq
echo -e "\n"

# 4. POST /users/{id}/tasks
curl -s -X POST "$BASE_URL/users/1/tasks" \
  -H "Content-Type: application/json" \
  -d '{
    "id": 10,
    "title": "Task Baru",
    "description": "Task dari nested route",
    "status": "progress",
    "deadline": "2025-12-31"
  }' | jq
echo -e "\n"

# 5. PATCH /tasks/{id}
curl -s -X PATCH "$BASE_URL/tasks/1" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }' | jq
echo -e "\n"

# 6. DELETE /tasks/{id}
curl -s -X DELETE "$BASE_URL/tasks/1" | jq
echo -e "\n"

================== USER ==================
# 7. GET /users
curl -s -X GET "$BASE_URL/users" | jq
echo -e "\n"

# 8. POST /users
curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "daz",
    "role": "admin"
  }' | jq
echo -e "\n"

# 9. DELETE /users/{id}
curl -s -X DELETE "$BASE_URL/users/1" | jq
echo -e "\n"
================== END ==================
