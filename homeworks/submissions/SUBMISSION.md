# 📝 Homework Submission - Day 2

**Họ tên:** Vũ Đặng Hải Đăng 

---

## ✅ Các bài đã hoàn thành

- [x] Bài 1: Statistics APIs (20 điểm)
- [x] Bài 2: Batch Create (25 điểm)
- [x] Bài 3: Batch Delete (20 điểm)
- [x] Bài 4: Connection Retry (25 điểm)
- [x] Bài 5: Health Check (15 điểm)
- [ ] Bài 6: Pagination (15 điểm - Bonus)
- [ ] Bài 7: Search (10 điểm - Bonus)

**Tổng điểm:** 105/105 (+ bonus)

---

## 🚀 Chạy project

```bash
# 1. Start Database
docker-compose up -d

# 2. Run Server
go run cmd/server/main.go

# 3. Test API (terminal khác)
curl http://localhost:8080/health 
```

---

## 📋 BÀI 1: STATISTICS APIs (20 điểm)

**Yêu cầu:** 
- Endpoint `GET /assets/stats` - Lấy thống kê tổng quát
- Endpoint `GET /assets/count` - Đếm assets với filters

### Test Instructions

#### Test 1.1: GET /assets/stats

```bash
curl http://localhost:8080/assets/stats
```

**Screenshot:**
![Test 1_1] (img/Test1_1.png)

---

#### Test 1.2: GET /assets/count (không filter)

```bash
curl http://localhost:8080/assets/count
```

**Screenshot:**
![Test 1_2] (img/Test1_2.png)

---

#### Test 1.3: GET /assets/count?type=domain

```bash
curl "http://localhost:8080/assets/count?type=domain"
```


**Screenshot:**
![Test 1_3] (img/Test1_3.png)

---

#### Test 1.4: GET /assets/count?type=domain&status=active

```bash
curl "http://localhost:8080/assets/count?type=domain&status=active"
```

**Screenshot:**
![Test 1_4] (img/Test1_4.png)

---


## 📋 BÀI 2: BATCH CREATE ASSETS (25 điểm)

**Yêu cầu:**
- Endpoint `POST /assets/batch` - Tạo nhiều assets trong 1 transaction
- All-or-nothing (rollback nếu có lỗi)
- Max 100 assets/request

### Test Instructions

#### Test 2.1: Batch Create Success (3 assets)

```bash
curl -X POST http://localhost:8080/assets/batch \
  -H "Content-Type: application/json" \
  -d '{
    "assets": [
      {"name":"batch1.com","type":"domain"},
      {"name":"batch2.com","type":"domain"},
      {"name":"192.168.1.1","type":"ip"}
    ]
  }'
```

**Screenshot:**
![Test 2_1] (img/Test2_1.png)

---

#### Test 2.2: Verify Assets Created

```bash
curl http://localhost:8080/assets | jq '.' 
```

**Screenshot:**
![Test 2_2] (img/Test2_2.png)

---

#### Test 2.3: Batch Create with Invalid Type (Rollback)

```bash
curl -X POST http://localhost:8080/assets/batch \
  -H "Content-Type: application/json" \
  -d '{
    "assets": [
      {"name":"good.com","type":"domain"},
      {"name":"bad.com","type":"invalid_type"}
    ]
  }' | jq
```

**Screenshot:**
![Test 2_3] (img/Test2_3.png)

---

## 📋 BÀI 3: BATCH DELETE ASSETS (20 điểm)

**Yêu cầu:**
- Endpoint `DELETE /assets/batch?ids=uuid1,uuid2,uuid3`
- Graceful handling (bỏ qua IDs không tồn tại)

### Test Instructions

#### Test 3.1: Create Test Data

```bash
ID1=$(curl -s -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{"name":"delete1.com","type":"domain"}' | jq -r '.id')

```

#### Test 3.2: Batch Delete (1 real + 1 fake ID)

```bash
curl -X DELETE "http://localhost:8080/assets/batch?ids=$ID1,fake-uuid-123"
```
**Screenshot:**
![Test 3_2] (img/Test3_2.png)

---

#### Test 3.3: Verify Deletion

```bash
curl http://localhost:8080/assets/$ID1
```

**Screenshot:**
![Test 3_3] (img/Test3_3.png)

---

## 📋 BÀI 4: DATABASE CONNECTION RETRY (25 điểm)

**Yêu cầu:**
- Auto retry khi DB connection thất bại
- Max 5 lần retry
- Exponential backoff: 1s → 2s → 4s → 8s → 16s

### Test Instructions

#### Test 4.1: Normal Startup (Quick)

```bash
go run cmd/server/main.go
```

**Screenshot:**
![Test 4_1] (img/Test4_1.png)

---

#### Test 4.2: Retry on Database Down

Dừng db docker và chạy lại trong lúc đang retry

**Screenshot:**
![Test 4_2] (img/Test4_2.png)

---


## 📋 BÀI 5: DATABASE HEALTH CHECK (15 điểm)

**Yêu cầu:**
- Enhanced `/health` endpoint
- Include database status + connection pool stats
- HTTP 200 when connected, 503 when down

### Test Instructions

#### Test 5.1: Health Check (Database Connected)

```bash
curl http://localhost:8080/health 
```

**Screenshot:**
![Test 5_1] (img/Test5_1.png)

---

#### Test 5.2: Health Check (Database Down)

```bash
docker-compose stop db
sleep 2
curl http://localhost:8080/health
```

**Screenshot:**
![Test 5_2] (img/Test5_2.png)

---

#### Test 5.3: Health Check Recovery

```bash
docker-compose start db
sleep 2
curl http://localhost:8080/health
```

**Screenshot:**
![Test 5_3] (img/Test5_3.png)

---
## 📊 Summary

| Bài | Điểm | Status | Ghi chú |
|-----|------|--------|--------|
| Bài 1 | 20 | ✅ | Statistics APIs |
| Bài 2 | 25 | ✅ | Batch Create + Rollback |
| Bài 3 | 20 | ✅ | Batch Delete |
| Bài 4 | 25 | ✅ | Retry + Exponential Backoff |
| Bài 5 | 15 | ✅ | Health Check |
| **Tổng** | **105** | | |

---

## 📝 Ghi chú


---

**Ready for review by:** dinhmanhtan (dmtangtnd@gmail.com)