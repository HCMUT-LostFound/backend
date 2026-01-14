# 🔧 Sửa CORS cho Mobile App

## ❌ Vấn Đề

Backend đang chặn requests từ mobile app vì:
1. CORS chỉ cho phép một số origins cụ thể
2. Mobile app (React Native/Expo) không có origin cụ thể
3. HEAD requests trả về 404 (preflight CORS)

## ✅ Giải Pháp

Đã cập nhật CORS config để:
1. Cho phép tất cả origins khi set `ALLOW_ALL_ORIGINS=true`
2. Hỗ trợ multiple origins từ `CORS_ALLOWED_ORIGINS` (phân cách bằng dấu phẩy)
3. Log CORS config để debug

## 🚀 Cách Sử Dụng

### Option 1: Cho phép tất cả origins (Khuyến nghị cho mobile app)

Trong Render Dashboard, thêm environment variable:
```
ALLOW_ALL_ORIGINS=true
```

### Option 2: Chỉ cho phép specific origins

Trong Render Dashboard, set:
```
CORS_ALLOWED_ORIGINS=https://your-web-app.com,https://another-domain.com
```

**Lưu ý:** Có thể set cả hai, nhưng `ALLOW_ALL_ORIGINS=true` sẽ override `CORS_ALLOWED_ORIGINS`.

## 📋 Cập Nhật Render Environment Variables

1. Vào Render Dashboard
2. Chọn service `lostfound-api`
3. Vào tab **Environment**
4. Thêm hoặc cập nhật:
   - `ALLOW_ALL_ORIGINS` = `true` (cho mobile app)
   - Hoặc `CORS_ALLOWED_ORIGINS` = `origin1,origin2` (cho web app)

## 🔄 Deploy Lại

Sau khi cập nhật environment variables, Render sẽ tự động redeploy. Hoặc bạn có thể:
1. Click **Manual Deploy** > **Deploy latest commit**
2. Hoặc push code mới lên repo

## ✅ Kiểm Tra

Sau khi deploy, test lại:
```bash
# Test từ mobile app
# Hoặc test bằng curl:
curl -X GET https://backend-lxgx.onrender.com/api/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Nếu vẫn lỗi, kiểm tra logs trong Render Dashboard.

---

**Sau khi deploy, test lại mobile app!**

