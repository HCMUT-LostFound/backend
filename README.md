# 📦 Lost & Found Backend API

Backend service cho hệ thống **Lost & Found App**, xây dựng bằng **Go (Gin Framework)**, sử dụng **PostgreSQL** làm cơ sở dữ liệu và hỗ trợ chạy bằng **Docker Compose**.

Dự án cung cấp các RESTful API để:
- Quản lý người dùng
- Đăng và lấy danh sách đồ thất lạc / tìm thấy (items)
- Xác thực người dùng (tích hợp bên ngoài như Clerk ở tầng middleware nếu có)

---

## 🚀 Công nghệ sử dụng

- Go 1.22+
- Gin (HTTP web framework)
- PostgreSQL 16
- Docker & Docker Compose
- SQL Migrations

---
## ▶️ Chạy dự án

### 1️⃣ Chạy PostgreSQL bằng Docker

    docker compose up -d

### 2️⃣ Cài dependencies

    go mod tidy

### 3️⃣ Chạy migration

    migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/lostfound?sslmode=disable" up

### 4️⃣ Chạy server

    go run cmd/api/main.go
