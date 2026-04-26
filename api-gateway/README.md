# API Gateway - Tugas Besar RPL 2

## Kelompok 3 - API Gateway / Integrator

Middleware untuk komunikasi antar-aplikasi dalam ekosistem UMKM.

### 🏗️ Arsitektur

```
User ──→ [Marketplace / POS / SupplierHub]
              │
              ▼
        [API Gateway]  ← Routing, Validasi JWT, Logging, Fee 0.5%
              │
              ▼
         [SmartBank]   ← Proses pembayaran, pajak/fee
              │
         ┌────┴────┐
         ▼         ▼
   [LogistiKita]  [UMKM Insight]
```

### 📋 Fitur Utama

| Fitur | Deskripsi |
|:------|:----------|
| **API Routing** | Meneruskan request ke service tujuan (SmartBank, Marketplace, POS, dll) |
| **Validasi JWT** | Memvalidasi token JWT pada setiap request |
| **Logging** | Mencatat semua request/response ke database |
| **Gateway Fee** | Memotong 0.5% dari setiap transaksi |
| **Rate Limiting** | Cooldown 10-30 detik & max 10 transaksi/hari |
| **Service Registry** | Mengelola daftar service yang terdaftar |
| **Health Check** | Memeriksa status kesehatan semua service |
| **Dashboard** | Statistik gateway (total request, fee, response time) |

### 🛠️ Tech Stack

- **Bahasa:** Go (Golang)
- **Framework:** Gin
- **Database:** PostgreSQL + GORM
- **Auth:** JWT (golang-jwt)
- **Pattern:** MVC / Clean Architecture

### 📁 Struktur Project

```
api-gateway/
├── config/          # Konfigurasi app & database
├── controller/      # Handler HTTP (MVC - Controller)
├── middleware/       # JWT, Logging, Rate Limit, CORS
├── model/           # Struct database & DTO (MVC - Model)
├── repository/      # Akses database (Repository Pattern)
├── routes/          # Definisi routing
├── service/         # Business logic (Service Layer)
├── utils/           # Helper functions
├── logs/            # Log files
├── .env             # Environment variables
├── go.mod           # Go module
└── main.go          # Entry point
```

### 🚀 Cara Menjalankan

1. **Setup Database PostgreSQL:**
   ```sql
   CREATE DATABASE api_gateway;
   ```

2. **Konfigurasi .env:** Edit file `.env` sesuai environment kamu.

3. **Install dependencies:**
   ```bash
   go mod tidy
   ```

4. **Jalankan:**
   ```bash
   go run main.go
   ```

5. **Server berjalan di:** `http://localhost:8080`

### 📡 API Endpoints

#### Public (Tanpa JWT)
| Method | Endpoint | Deskripsi |
|:-------|:---------|:----------|
| GET | `/health` | Health check gateway |
| POST | `/auth/token` | Generate JWT token |

#### Gateway Management (JWT Required)
| Method | Endpoint | Deskripsi |
|:-------|:---------|:----------|
| GET | `/gateway/dashboard` | Dashboard statistik |
| GET | `/gateway/stats` | Statistik request |
| GET | `/gateway/logs` | Lihat log request |
| GET | `/gateway/logs/:request_id` | Detail log |
| POST | `/gateway/fee/calculate` | Hitung fee 0.5% |
| GET | `/gateway/fee/stats` | Statistik fee |
| GET | `/gateway/fee/:transaction_id` | Detail fee |
| GET | `/gateway/fees` | Semua fee |
| PUT | `/gateway/fee/:id/status` | Update status fee |
| GET | `/gateway/services` | Daftar service |
| POST | `/gateway/services` | Daftarkan service baru |
| PUT | `/gateway/services/:name/status` | Update status service |
| DELETE | `/gateway/services/:name` | Hapus service |
| GET | `/gateway/health` | Health check semua service |
| GET | `/gateway/health/:name` | Health check satu service |

#### Proxy Routes (JWT + Rate Limit)
| Method | Endpoint | Deskripsi |
|:-------|:---------|:----------|
| ANY | `/api/smartbank/*` | Forward ke SmartBank |
| ANY | `/api/marketplace/*` | Forward ke Marketplace |
| ANY | `/api/pos/*` | Forward ke POS |
| ANY | `/api/logistikita/*` | Forward ke LogistiKita |
| ANY | `/api/supplierhub/*` | Forward ke SupplierHub |
| ANY | `/api/umkm-insight/*` | Forward ke UMKM Insight |

### 📝 Contoh Penggunaan

#### 1. Generate Token
```bash
curl -X POST http://localhost:8080/auth/token \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "username": "richard", "role": "user"}'
```

#### 2. Proxy Request ke SmartBank
```bash
curl -X GET http://localhost:8080/api/smartbank/saldo \
  -H "Authorization: Bearer <TOKEN>"
```

#### 3. Hitung Fee Gateway
```bash
curl -X POST http://localhost:8080/gateway/fee/calculate \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"amount": 100000, "source_service": "marketplace", "destination_service": "smartbank", "user_id": 1}'
```

### 💰 Aturan Keuangan Gateway

| Parameter | Nilai |
|:----------|:------|
| Fee Gateway | 0.5% per transaksi |
| Cooldown Transaksi | 10-30 detik |
| Max Transaksi Harian | 10 per user |
