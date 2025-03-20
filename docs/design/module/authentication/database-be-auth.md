# 🔐 Operasi Database pada Autentikasi Backend - Dokumentasi

## 📑 Daftar Isi
- [🔐 Operasi Database pada Autentikasi Backend - Dokumentasi](#-operasi-database-pada-autentikasi-backend---dokumentasi)
  - [📑 Daftar Isi](#-daftar-isi)
  - [📝 Pendahuluan](#-pendahuluan)
  - [💾 Struktur Database](#-struktur-database)
  - [📊 Daftar Operasi Database](#-daftar-operasi-database)
    - [1️⃣ Pemeriksaan Email](#1️⃣-pemeriksaan-email)
    - [2️⃣ Penyimpanan Token Aktivasi](#2️⃣-penyimpanan-token-aktivasi)
    - [3️⃣ Validasi Token Aktivasi](#3️⃣-validasi-token-aktivasi)
    - [4️⃣ Pembuatan User Baru](#4️⃣-pembuatan-user-baru)
    - [5️⃣ Pencarian User](#5️⃣-pencarian-user)
    - [6️⃣ Pembaruan Password](#6️⃣-pembaruan-password)
  - [🔒 Prinsip ACID dalam Transaksi Database](#-prinsip-acid-dalam-transaksi-database)
  - [⚠️ Penanganan Error Database](#️-penanganan-error-database)
  - [📘 Istilah Teknis](#-istilah-teknis)

## 📝 Pendahuluan

Dokumen ini menjelaskan semua operasi database yang digunakan dalam sistem autentikasi backend. Operasi-operasi ini menjamin integritas dan keamanan data melalui penerapan prinsip ACID (Atomicity, Consistency, Isolation, Durability).

## 💾 Struktur Database

```sql
-- Tabel users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    fullname VARCHAR(30) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel activation_tokens
CREATE TABLE activation_tokens (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL, -- 'registration' atau 'password_reset'
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 📊 Daftar Operasi Database

Berikut adalah daftar semua operasi database yang digunakan dalam sistem autentikasi:

### 1️⃣ Pemeriksaan Email

**Nama Fungsi**: `database.CheckExistingEmail(email)`

**Fungsi**: Memeriksa apakah email sudah terdaftar dalam database

**Parameter**:
- `email` (string): Alamat email yang ingin diperiksa

**Query SQL**:
```sql
SELECT 1 FROM users WHERE email = $1
```

**Return Type**:
- *customerror.CustomError

**Error yang Mungkin Muncul**:
- ❌ Postgres Error (kode: 500) - Database error
- ❌ Custom Error (kode: sesuai kasus) - Error spesifik

**Digunakan pada**:
- Saat pendaftaran untuk mencegah duplikasi email
- Saat reset password untuk memastikan email terdaftar

---

### 2️⃣ Penyimpanan Token Aktivasi

**Nama Fungsi**: `authPostgres.SaveActivationToken(req *ActivationTokenRequest)`

**Fungsi**: Menyimpan token aktivasi untuk verifikasi email

**Parameter Type**:
```go
type ActivationTokenRequest struct {
    Email         string
    TokenHash     string
    TokenType     string
    ExpiryMinutes int
}
```

**Query SQL**:
```sql
INSERT INTO activation_tokens (email, token_hash, type, expires_at) 
VALUES ($1, $2, $3, NOW() + INTERVAL '$4 minutes')
```

**Return Type**:
- *customerror.CustomError

**Error yang Mungkin Muncul**:
- ❌ Postgres Error - Database error (kode: sesuai dengan customerror.NewPostgresError)

---

### 3️⃣ Validasi Token Aktivasi

**Nama Fungsi**: `database.ValidateActivationToken(req *ActivationTokenRequest)`

**Parameter Type**:
```go
type ActivationTokenRequest struct {
    Email     string
    TokenHash string
    TokenType string
}
```

**Query SQL**:
```sql
SELECT id FROM activation_tokens 
WHERE email = $1 AND token_hash = $2 AND type = $3 AND expires_at > NOW()
```

**Return Type**:
- *customerror.CustomError

**Error yang Mungkin Muncul**:
- ❌ "token not found or expired" (kode: 404)
- ❌ Postgres Error (kode: 500) - Database error

---

### 4️⃣ Pembuatan User Baru

**Nama Fungsi**: `database.CreateUser(req *CreateUserRequest)`

**Parameter Type**:
```go
type CreateUserRequest struct {
    Email        string
    Fullname     string
    PasswordHash string
}
```

**Query SQL**:
```sql
-- In transaction:
-- Query 1: Insert data pengguna
INSERT INTO users (email, fullname, password_hash) VALUES ($1, $2, $3)

-- Query 2: Hapus token aktivasi
DELETE FROM activation_tokens WHERE email = $1
```

**Return Type**:
- *customerror.CustomError

**Properti ACID**:
- **Atomic**: Kedua operasi dalam satu transaksi
- **Consistency**: Menjamin data sesuai dengan aturan validasi
- **Isolation**: Level ReadCommitted
- **Durability**: Perubahan disimpan permanen setelah commit

**Error yang Mungkin Muncul**:
- ❌ Postgres Error - Database error (kode: sesuai dengan customerror.NewPostgresError)

---

### 5️⃣ Pencarian User

**Nama Fungsi**: `database.GetUserByEmail(email)`

**Parameter**:
- `email` (string): Alamat email pengguna

**Return Type**:
```go
type User struct {
    ID           string
    Email        string
    Fullname     string
    PasswordHash string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

**Query SQL**:
```sql
SELECT id, email, fullname, password_hash, created_at, updated_at 
FROM users WHERE email = $1
```

**Error yang Mungkin Muncul**:
- ❌ "user not found" (kode: 404)
- ❌ Postgres Error - Database error (kode: sesuai dengan customerror.NewPostgresError)

---

### 6️⃣ Pembaruan Password

**Nama Fungsi**: `database.UpdateUserPassword(req *UpdatePasswordRequest)`

**Parameter Type**:
```go
type UpdatePasswordRequest struct {
    Email           string
    NewPasswordHash string
}
```

**Query SQL**:
```sql
-- In transaction:
-- Query 1: Update password
UPDATE users 
SET password_hash = $1, updated_at = CURRENT_TIMESTAMP 
WHERE email = $2

-- Query 2: Hapus token aktivasi
DELETE FROM activation_tokens WHERE email = $1
```

**Return Type**:
- *customerror.CustomError

**Error yang Mungkin Muncul**:
- ❌ "user not found" (kode: 404)
- ❌ Postgres Error (kode: 500) - Database error

## 🔒 Prinsip ACID dalam Transaksi Database

Sistem autentikasi menerapkan 4 prinsip utama dalam transaksi database:

1. **Atomicity (Keutuhan)**
   - Transaksi dijalankan sepenuhnya atau tidak sama sekali
   - Jika ada error, rollback otomatis dilakukan untuk mengembalikan ke keadaan awal
   - Contoh: Saat membuat user, jika gagal menyimpan data, token tidak akan dihapus

2. **Consistency (Konsistensi)**
   - Data selalu dalam keadaan valid sesuai aturan bisnis
   - Constraint seperti UNIQUE dan NOT NULL selalu dipatuhi
   - Referential integrity terjaga (misalnya: relasi antar tabel)

3. **Isolation (Isolasi)**
   - Transaksi dijalankan secara terisolasi dari transaksi lain
   - Menggunakan level ReadCommitted untuk menghindari dirty read
   - Transaksi berjalan seolah-olah hanya ada satu transaksi dalam database

4. **Durability (Ketahanan)**
   - Data yang sudah dicommit tetap tersimpan meskipun sistem crash
   - Menggunakan write-ahead logging untuk menjamin ketahanan data
   - Perubahan disimpan secara permanen ke disk setelah commit

## ⚠️ Penanganan Error Database

Sistem menggunakan package customerror untuk menangani error. Berikut adalah daftar error umum:

| Error | Kode HTTP | Penanganan |
|-------|-----------|------------|
| Record not found | 404 | Kembalikan custom error dengan pesan spesifik |
| Database error | 500 | Gunakan NewPostgresError untuk error database |
| Validation error | 400 | Gunakan NewCustomError dengan pesan spesifik |
| Unique violation | 409 | Tangani melalui NewPostgresError |
| Foreign key violation | 400 | Tangani melalui NewPostgresError |

## 📘 Istilah Teknis

- **Transaksi**: Serangkaian operasi database yang diperlakukan sebagai satu kesatuan
- **Commit**: Menyimpan perubahan secara permanen ke database
- **Rollback**: Membatalkan perubahan dan kembali ke keadaan sebelumnya
- **Isolation Level**: Tingkat isolasi transaksi dari transaksi lain
- **ReadCommitted**: Level isolasi yang hanya membaca data yang sudah di-commit
- **Constraint**: Aturan untuk menjaga integritas data
- **Foreign Key**: Hubungan referensial antara dua tabel