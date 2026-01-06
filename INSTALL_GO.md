# Cài Đặt Go trên Windows

## Cách 1: Tải và Cài Đặt Thủ Công (Khuyến nghị)

### Bước 1: Tải Go

1. Truy cập: https://go.dev/dl/
2. Tải file: **`go1.25.5.windows-amd64.msi`** (hoặc phiên bản mới nhất)

### Bước 2: Cài Đặt

1. Chạy file `.msi` vừa tải
2. Làm theo hướng dẫn trong installer
3. Go sẽ tự động được thêm vào PATH

### Bước 3: Kiểm Tra

Mở PowerShell mới và chạy:

```powershell
go version
```

Bạn sẽ thấy: `go version go1.25.5 windows/amd64` (hoặc phiên bản tương ứng)

---

## Cách 2: Sử dụng Chocolatey (Cần quyền Admin)

Mở PowerShell với quyền Administrator và chạy:

```powershell
choco install golang -y
```

Sau đó khởi động lại PowerShell.

---

## Cách 3: Sử dụng Scoop

```powershell
scoop install go
```

---

## Sau Khi Cài Đặt

1. **Khởi động lại PowerShell** để PATH được cập nhật
2. Kiểm tra: `go version`
3. Chạy backend:

```powershell
cd C:\project\mobile\backend
go mod download
go run cmd/api/main.go
```

---

## Link Tải Trực Tiếp

**Latest Version:**
- https://go.dev/dl/go1.25.5.windows-amd64.msi

**Tất cả versions:** https://go.dev/dl/

