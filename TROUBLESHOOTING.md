# Hướng dẫn khắc phục lỗi DocuMind

## Lỗi thường gặp và cách khắc phục

### 1. Lỗi API Quota (Error 429)

**Triệu chứng:**
- Thông báo: "API quota đã hết. Vui lòng thử lại sau hoặc liên hệ admin để nâng cấp quota."
- Lỗi 429 từ Google Gemini API

**Nguyên nhân:**
- Đã vượt quá giới hạn sử dụng Google Gemini API
- Quota miễn phí đã hết (thường 15 requests/phút cho Gemini Pro)

**Cách khắc phục:**

#### Cho Admin:
1. **Kiểm tra quota hiện tại:**
   - Truy cập [Google AI Studio](https://makersuite.google.com/app/apikey)
   - Xem quota usage trong dashboard

2. **Nâng cấp quota:**
   - Nâng cấp lên Google AI Studio Pro ($20/tháng)
   - Hoặc tạo API key mới với billing account

3. **Cấu hình billing:**
   ```bash
   # Thêm billing account vào project
   gcloud billing projects link YOUR_PROJECT_ID --billing-account=YOUR_BILLING_ACCOUNT_ID
   ```

4. **Thay đổi API key:**
   - Tạo API key mới trong Google AI Studio
   - Cập nhật `GEMINI_API_KEY` trong file `.env`

#### Cho User:
- Thử lại sau vài phút
- Liên hệ admin để nâng cấp quota
- Sử dụng file nhỏ hơn để giảm token usage

### 2. Lỗi Xác thực API

**Triệu chứng:**
- "Lỗi xác thực API. Vui lòng kiểm tra cấu hình."

**Cách khắc phục:**
1. Kiểm tra `GEMINI_API_KEY` trong file `.env`
2. Đảm bảo API key hợp lệ và có quyền truy cập
3. Tạo API key mới nếu cần

### 3. Lỗi Định dạng File

**Triệu chứng:**
- "Định dạng file không được hỗ trợ. Vui lòng sử dụng PDF, DOC hoặc DOCX."

**Cách khắc phục:**
- Chỉ sử dụng file PDF, DOC, DOCX
- Kiểm tra file không bị hỏng
- Chuyển đổi file sang định dạng được hỗ trợ

### 4. Lỗi Kết nối Mạng

**Triệu chứng:**
- "Lỗi kết nối mạng. Vui lòng kiểm tra kết nối internet và thử lại."

**Cách khắc phục:**
- Kiểm tra kết nối internet
- Kiểm tra firewall/antivirus
- Thử lại sau vài phút

### 5. Lỗi Database

**Triệu chứng:**
- "Database query error" hoặc "Failed to save analysis"

**Cách khắc phục:**
1. Kiểm tra kết nối PostgreSQL
2. Kiểm tra `DATABASE_URL` trong `.env`
3. Restart database service
4. Kiểm tra disk space

## Cấu hình Môi trường

### File .env (Backend)
```env
GEMINI_API_KEY=your_google_gemini_api_key_here
DATABASE_URL=host=localhost user=postgres password=yourpass dbname=documind_db port=5432 sslmode=disable
PORT=8080
```

### Kiểm tra API Key
```bash
# Test API key
curl -H "Content-Type: application/json" \
     -d '{"contents":[{"parts":[{"text":"Hello"}]}]}' \
     "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=YOUR_API_KEY"
```

## Monitoring và Logs

### Backend Logs
```bash
# Xem logs real-time
tail -f backend/logs/app.log

# Tìm lỗi quota
grep -i "quota\|429\|exceeded" backend/logs/app.log
```

### Frontend Logs
- Mở Developer Tools (F12)
- Xem Console tab để debug lỗi

## Liên hệ Hỗ trợ

Nếu gặp lỗi không thể khắc phục:
1. Chụp màn hình lỗi
2. Copy log lỗi từ console
3. Liên hệ admin với thông tin:
   - Mô tả lỗi
   - File đang upload (nếu có)
   - Thời gian xảy ra lỗi
   - Log lỗi

## Tips Tối ưu

1. **Giảm token usage:**
   - Sử dụng file nhỏ hơn
   - Tách file lớn thành nhiều phần
   - Sử dụng cache (file đã upload trước đó)

2. **Tăng hiệu suất:**
   - Sử dụng SSD cho database
   - Tăng RAM cho server
   - Sử dụng CDN cho static files

3. **Backup và Recovery:**
   - Backup database thường xuyên
   - Lưu trữ file quan trọng offline
   - Test restore procedure 