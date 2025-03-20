# 🔐 Dokumentasi API Autentikasi

## 📑 Table of Contents
- [🔐 Dokumentasi API Autentikasi](#-dokumentasi-api-autentikasi)
  - [📑 Table of Contents](#-table-of-contents)
  - [🌐 Overview](#-overview)
    - [📬 General Headers for All Requests](#-general-headers-for-all-requests)
  - [💾 Database Schema](#-database-schema)
  - [⚙️ Environment Variables](#️-environment-variables)
  - [🔄 Third-Party Services](#-third-party-services)
  - [📦 External Packages](#-external-packages)
  - [📋 Ringkasan Endpoint API](#-ringkasan-endpoint-api)
  - [🔌 API Endpoints](#-api-endpoints)
    - [1️⃣ Request Token](#1️⃣-request-token)
    - [2️⃣ Register User](#2️⃣-register-user)
    - [3️⃣ Login](#3️⃣-login)
    - [4️⃣ Reset Password](#4️⃣-reset-password)
    - [5️⃣ Refresh Token](#5️⃣-refresh-token)
  - [⚠️ Common Error Responses](#️-common-error-responses)

## 🌐 Overview

Modul Autentikasi menyediakan endpoint-endpoint untuk manajemen akses pengguna, termasuk registrasi, login, logout, dan reset password. Modul ini menggunakan JWT untuk autentikasi dan mengintegrasikan layanan email untuk pengiriman token verifikasi.

### 📬 General Headers for All Requests

| Header | Format | Keterangan |
|--------|--------|------------|
| **Content-Type** | `application/json` | Wajib untuk semua request |
| **Authorization** | `Bearer {accessToken}` | Diperlukan untuk endpoint yang membutuhkan autentikasi |

## 💾 Database Schema

```sql
-- Tabel users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(30) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel activation_tokens
CREATE TABLE activation_tokens (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(6) NOT NULL,
    type VARCHAR(20) NOT NULL, -- 'registration' atau 'password_reset'
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## ⚙️ Environment Variables

```
# Server
PORT=8080
ENV=development
JWT_SECRET=your_jwt_secret
ACCESS_TOKEN_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=7d
ACTIVATION_TOKEN_EXPIRY=15m

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=auth_db

# Email Service
RESEND_API_KEY=your_resend_api_key
RESEND_FROM_EMAIL=noreply@yourdomain.com
```

## 🔄 Third-Party Services

| Layanan | Kegunaan | Environment Variables |
|---------|----------|----------------------|
| Resend | Pengiriman email verifikasi dan reset password | `RESEND_API_KEY`, `RESEND_FROM_EMAIL` |

## 📦 External Packages

| Package | Kegunaan |
|---------|----------|
| github.com/gin-gonic/gin | Web framework |
| github.com/go-pg/pg/v10 | PostgreSQL ORM | // using raw SQL
| github.com/golang-jwt/jwt | JWT authentication |
| github.com/resend/resend-go | Resend email client |
| golang.org/x/crypto/bcrypt | Password hashing |

## 📋 Ringkasan Endpoint API

| Endpoint | Metode | Deskripsi | Auth yang Diperlukan | Potensi Error |
|----------|--------|-----------|---------------------|---------------|
| `/auth/token/:type` | POST | Request token (registrasi/reset password) | Tidak | DB, Resend (3rd) |
| `/auth/register` | POST | Registrasi user baru | Tidak | DB |
| `/auth/login` | POST | Login user | Tidak | DB |
| `/auth/forget-password` | POST | Reset password | Tidak | DB |
| `/auth/refresh-token` | POST | Refresh access token | Ya (refresh token) | - |
## 🔌 API Endpoints

### 1️⃣ Request Token

> Endpoint untuk meminta token aktivasi yang akan dikirim ke email pengguna untuk proses registrasi atau reset password

**🔹 Endpoint:** `POST /auth/token/:type`

**🔹 URL Parameters:**
- type: 'registration' atau 'forget-password'

**🔹 Request Body:**
```json
{
  "email": "user@example.com"
}
```

**🔹 Success Response (200):**
```json
{
  "status": 200,
  "message": "Kode aktivasi telah dikirim ke email"
}
```

**🔹 Error Responses:**

<details>
<summary><strong>Bad Request (400)</strong></summary>

```json
{
  "status": 400,
  "message": "Format email tidak valid"
}
```

```json
{
  "status": 400,
  "message": "Tipe token tidak valid"
}
```
</details>

<details>
<summary><strong>Conflict (409)</strong></summary>

```json
{
  "status": 409,
  "message": "Email sudah terdaftar"
}
```
</details>

<details>
<summary><strong>Not Found (404)</strong></summary>

```json
{
  "status": 404,
  "message": "Email tidak terdaftar"
}
```
</details>

**🔹 Packages:**
- github.com/gin-gonic/gin
- github.com/go-pg/pg/v10

**🔹 Services:**
- Resend (3rd) - Email delivery service

**🔹 Handler Operations:**

* GetParam("type") - Validasi tipe token
  * jika bukan "registration" atau "forget-password" respon `{"status": 400, "message": "Tipe token tidak valid"}`
* service.ValidateEmail(email) - Validasi format email
  * jika tidak valid respon `{"status": 400, "message": "Format email tidak valid"}`
* Jika type == "registration":
  * database.CheckExistingEmail(email) - Memeriksa email yang sudah terdaftar
    * ACID: Consistency, Durability dengan Isolation-ReadCommitted
    * Query: `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
    * jika email sudah ada respon `{"status": 409, "message": "Email sudah terdaftar"}`
    * jika database tidak tersedia respon `{"status": 503, "message": "Layanan sementara tidak tersedia"}`
    * jika koneksi database timeout respon `{"status": 504, "message": "Waktu koneksi ke database habis"}`
* Jika type == "forget-password":
  * database.CheckExistingEmail(email) - Memeriksa email yang sudah terdaftar
    * ACID: Consistency, Durability dengan Isolation-ReadCommitted
    * Query: `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
    * jika email tidak ada respon `{"status": 404, "message": "Email tidak terdaftar"}`
    * jika database tidak tersedia respon `{"status": 503, "message": "Layanan sementara tidak tersedia"}`
    * jika koneksi database timeout respon `{"status": 504, "message": "Waktu koneksi ke database habis"}`
* service.GenerateActivationToken() - Generate token 6 digit
* database.SaveActivationToken(email, token, type) - Menyimpan token aktivasi
  * Membuat transaksi database (ACID: Atomic, Consistency, Isolation-ReadCommitted, Durability)
  * Query: `INSERT INTO activation_tokens (email, token, type, expires_at) VALUES ($1, $2, $3, NOW() + INTERVAL '15 minutes')`
  * Commit transaksi
  * jika gagal: rollback dan respon `{"status": 500, "message": "Gagal menyimpan token aktivasi"}`
  * jika database tidak tersedia respon `{"status": 503, "message": "Layanan sementara tidak tersedia"}`
  * jika koneksi database timeout respon `{"status": 504, "message": "Waktu koneksi ke database habis"}`
* Resend.SendActivationEmail(email, token, type) - Mengirim email dengan token
  * jika gagal mengirim email respon `{"status": 500, "message": "Gagal mengirim email aktivasi"}`
  * jika layanan email tidak tersedia respon `{"status": 502, "message": "Gagal menghubungi layanan email"}`
  * jika timeout respon `{"status": 504, "message": "Waktu koneksi ke layanan email habis"}`
* Response sukses: `{"status": 200, "message": "Kode aktivasi telah dikirim ke email"}`

### 2️⃣ Register User

> Endpoint untuk mendaftarkan pengguna baru dengan token aktivasi

**🔹 Endpoint:** `POST /auth/register`

**🔹 Request Body:**
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "password_confirmation": "password123",
  "activation_code": "123456"
}
```

**🔹 Success Response (201):**
```json
{
  "status": 201,
  "message": "Registrasi berhasil",
  "data": {
    "user_id": "uuid-string",
    "email": "user@example.com",
    "username": "username"
  }
}
```

**🔹 Error Responses:**

<details>
<summary><strong>Bad Request (400)</strong></summary>

```json
{
  "status": 400,
  "message": "Validasi gagal",
  "errors": {
    "username": "Username maksimal 30 karakter",
    "password": "Password tidak cocok dengan konfirmasi"
  }
}
```
</details>

<details>
<summary><strong>Not Found (404)</strong></summary>

```json
{
  "status": 404,
  "message": "Token aktivasi tidak valid atau sudah kadaluarsa"
}
```
</details>

**🔹 Packages:**
- github.com/gin-gonic/gin
- github.com/go-pg/pg/v10
- golang.org/x/crypto/bcrypt

**🔹 Services:**
- None

**🔹 Handler Operations:**

* service.ValidateRegistrationInput(request) - Validasi input registrasi
  * validasi panjang username <= 30 karakter
  * validasi kecocokan password
  * validasi format email
  * jika tidak valid respon `{"status": 400, "message": "Validasi gagal", "errors": {detail}}`
* database.ValidateActivationToken(email, token) - Validasi token aktivasi
  * ACID: Consistency, Durability dengan Isolation-ReadCommitted
  * Query: `SELECT id FROM activation_tokens WHERE email = $1 AND token = $2 AND type = 'registration' AND expires_at > NOW()`
  * jika tidak valid respon `{"status": 404, "message": "Token aktivasi tidak valid atau sudah kadaluarsa"}`
  * jika database tidak tersedia respon `{"status": 503, "message": "Layanan sementara tidak tersedia"}`
  * jika koneksi database timeout respon `{"status": 504, "message": "Waktu koneksi ke database habis"}`
* service.HashPassword(password) - Hash password
* database.CreateUser(user_data) - Menyimpan data pengguna
  * Membuat transaksi database (ACID: Atomic, Consistency, Isolation-ReadCommitted, Durability)
  * Query 1: `INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3)`
  * Query 2: `DELETE FROM activation_tokens WHERE email = $1`
  * Commit transaksi
  * jika gagal: rollback dan respon `{"status": 500, "message": "Gagal membuat user"}`
  * jika database tidak tersedia respon `{"status": 503, "message": "Layanan sementara tidak tersedia"}`
  * jika koneksi database timeout respon `{"status": 504, "message": "Waktu koneksi ke database habis"}`
* Response sukses: `{"status": 201, "message": "Registrasi berhasil", "data": user_data}`

### 3️⃣ Login

> Endpoint untuk melakukan autentikasi pengguna dan mendapatkan token akses

**🔹 Endpoint:** `POST /auth/login`

**🔹 Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**🔹 Success Response (200):**
```json
{
  "status": 200,
  "message": "Login berhasil",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

**🔹 Error Responses:**

<details>
<summary><strong>Unauthorized (401)</strong></summary>

```json
{
  "status": 401,
  "message": "Email atau password salah"
}
```
</details>

**🔹 Packages:**
- github.com/gin-gonic/gin
- github.com/go-pg/pg/v10
- golang.org/x/crypto/bcrypt
- github.com/golang-jwt/jwt

**🔹 Services:**
- None

**🔹 Handler Operations:**

* database.GetUserByEmail(email) - Mengambil data pengguna
  * ACID: Consistency, Durability dengan Isolation-ReadCommitted
  * Query: `SELECT id, email, username, password_hash FROM users WHERE email = $1`
  * jika user tidak ditemukan respon `{"status": 401, "message": "Email atau password salah"}`
  * jika database tidak tersedia respon `{"status": 503, "message": "Layanan sementara tidak tersedia"}`
  * jika koneksi database timeout respon `{"status": 504, "message": "Waktu koneksi ke database habis"}`
* service.VerifyPassword(password_hash, password) - Verifikasi password
  * jika tidak cocok respon `{"status": 401, "message": "Email atau password salah"}`
* service.GenerateTokenPair(user_id) - Generate access dan refresh token
  * Generate access token dengan expiry 15 menit
  * Generate refresh token dengan expiry 7 hari
* Response sukses: `{"status": 200, "message": "Login berhasil", "data": token_data}`

### 4️⃣ Reset Password

> Endpoint untuk melakukan reset password dengan token aktivasi

**🔹 Endpoint:** `POST /auth/forget-password`

**🔹 Request Body:**
```json
{
  "email": "user@example.com",
  "new_password": "newpassword123",
  "new_password_confirmation": "newpassword123",
  "activation_code": "123456"
}
```

**🔹 Success Response (200):**
```json
{
  "status": 200,
  "message": "Password berhasil diubah"
}
```

**🔹 Error Responses:**

<details>
<summary><strong>Bad Request (400)</strong></summary>

```json
{
  "status": 400,
  "message": "Password baru tidak cocok dengan konfirmasi"
}
```
</details>

<details>
<summary><strong>Not Found (404)</strong></summary>

```json
{
  "status": 404,
  "message": "Token aktivasi tidak valid atau sudah kadaluarsa"
}
```
</details>

**🔹 Packages:**
- github.com/gin-gonic/gin
- github.com/go-pg/pg/v10
- golang.org/x/crypto/bcrypt

**🔹 Services:**
- None

**🔹 Handler Operations:**

* service.ValidatePasswordInput(request) - Validasi input password baru
  * validasi kecocokan password dengan konfirmasi
  * jika tidak valid respon `{"status": 400, "message": "Password baru tidak cocok dengan konfirmasi"}`
* database.ValidateActivationToken(email, token) - Validasi token aktivasi
  * ACID: Consistency, Durability dengan Isolation-ReadCommitted
  * Query: `SELECT id FROM activation_tokens WHERE email = $1 AND token = $2 AND type = 'forget-password' AND expires_at > NOW()`
  * jika tidak valid respon `{"status": 404, "message": "Token aktivasi tidak valid atau sudah kadaluarsa"}`
  * jika database tidak tersedia respon `{"status": 503, "message": "Layanan sementara tidak tersedia"}`
  * jika koneksi database timeout respon `{"status": 504, "message": "Waktu koneksi ke database habis"}`
* service.HashPassword(new_password) - Hash password baru
* database.UpdateUserPassword(email, new_password_hash) - Update password user
  * Membuat transaksi database (ACID: Atomic, Consistency, Isolation-ReadCommitted, Durability)
  * Query 1: `UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE email = $2`
  * Query 2: `DELETE FROM activation_tokens WHERE email = $2`
  * Commit transaksi
  * jika gagal: rollback dan respon `{"status": 500, "message": "Gagal mengubah password"}`
  * jika database tidak tersedia respon `{"status": 503, "message": "Layanan sementara tidak tersedia"}`
  * jika koneksi database timeout respon `{"status": 504, "message": "Waktu koneksi ke database habis"}`
* Response sukses: `{"status": 200, "message": "Password berhasil diubah"}`

### 5️⃣ Refresh Token

> Endpoint untuk memperbarui access token menggunakan refresh token

**🔹 Endpoint:** `POST /auth/refresh-token`

**🔹 Packages:**
- github.com/gin-gonic/gin
- github.com/golang-jwt/jwt

**🔹 Services:**
- None

**🔹 Headers:**
- Cookie: refresh_token=eyJhbGciOiJIUzI1NiIs... (httpOnly)

**🔹 Success Response (200):**
```json
{
  "status": 200,
  "message": "Token berhasil diperbarui",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

**🔹 Error Responses:**

<details>
<summary><strong>Unauthorized (401)</strong></summary>

```json
{
  "status": 401,
  "message": "Refresh token tidak valid atau sudah kadaluarsa"
}
```
</details>

**🔹 Handler Operations:**

* service.ValidateRefreshToken(refresh_token) - Validasi refresh token dari cookie
  * jika tidak valid respon `{"status": 401, "message": "Refresh token tidak valid atau sudah kadaluarsa"}`
* service.GenerateAccessToken(user_id) - Generate access token baru
* Response sukses: `{"status": 200, "message": "Token berhasil diperbarui", "data": token_data}`

## ⚠️ Common Error Responses

Berikut adalah beberapa respons error umum yang mungkin terjadi pada semua endpoint:

| Kode | Deskripsi | Response Body | Penyebab |
|------|-----------|--------------|----------|
| **500** | Internal Server Error | `{"status": 500, "message": "Terjadi kesalahan pada server"}` | Error pemrosesan pada server |
| **503** | Service Unavailable | `{"status": 503, "message": "Layanan sementara tidak tersedia"}` | **Database tidak tersedia** |
| **502** | Bad Gateway | `{"status": 502, "message": "Kesalahan koneksi ke layanan eksternal"}` | **Resend (3rd) tidak tersedia** |
| **504** | Gateway Timeout | `{"status": 504, "message": "Waktu koneksi ke layanan eksternal habis"}` | **Timeout database atau Resend (3rd)** |
| **408** | Request Timeout | `{"status": 408, "message": "Waktu permintaan habis"}` | Waktu pemrosesan habis |
| **429** | Too Many Requests | `{"status": 429, "message": "Terlalu banyak permintaan, silakan coba lagi nanti"}` | Rate limiting diaktifkan |