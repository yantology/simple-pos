# 🔐 Layanan Autentikasi Backend - Dokumentasi Services

## 📑 Daftar Isi
- [🔐 Layanan Autentikasi Backend - Dokumentasi Services](#-layanan-autentikasi-backend---dokumentasi-services)
  - [📑 Daftar Isi](#-daftar-isi)
  - [📝 Pendahuluan](#-pendahuluan)
  - [🧰 Interface Layanan](#-interface-layanan)
  - [🛠️ Implementasi Layanan](#️-implementasi-layanan)
    - [1️⃣ Validasi Email](#1️⃣-validasi-email)
    - [2️⃣ Pembuat Token Aktivasi](#2️⃣-pembuat-token-aktivasi)
    - [3️⃣ Validasi Data Registrasi](#3️⃣-validasi-data-registrasi)
    - [4️⃣ Pengaman Password](#4️⃣-pengaman-password)
    - [5️⃣ Pemeriksa Password](#5️⃣-pemeriksa-password)
    - [6️⃣ Validasi Input Password](#6️⃣-validasi-input-password)
    - [7️⃣ Pembuat Token](#7️⃣-pembuat-token)
    - [8️⃣ Validasi Token Claims](#8️⃣-validasi-token-claims)
    - [9️⃣ Generate Access Token](#9️⃣-generate-access-token)
  - [⚠️ Kode Error Umum](#️-kode-error-umum)
  - [📘 Istilah Teknis](#-istilah-teknis)

## 📝 Pendahuluan

Dokumen ini menjelaskan semua layanan (services) yang digunakan dalam sistem autentikasi backend aplikasi. Services ini menangani berbagai proses autentikasi seperti registrasi, validasi, dan manajemen token.

## 🧰 Interface Layanan

Interface `AuthService` mendefinisikan semua operasi yang tersedia:

```go
type AuthService interface {
    // Email validation and token generation
    ValidateEmail(email string) *customerror.CustomError
    GenerateActivationToken() (string, *customerror.CustomError)
    
    // User registration and password management
    ValidateRegistrationInput(email, username, password, passwordConfirmation string) *customerror.CustomError
    HashPassword(password string) (string, *customerror.CustomError)
    VerifyPassword(hashedPassword, password string) *customerror.CustomError
    ValidatePasswordInput(password, passwordConfirmation string) *customerror.CustomError

    // Token operations
    GenerateTokenPair(userID, email, userType string) (accessToken, refreshToken string, err *customerror.CustomError)
    ValidateTokenClaims(token string) (*jwtPkg.TokenClaims, *customerror.CustomError)
    GenerateAccessToken(userID, email, userType string) (string, *customerror.CustomError)
}
```

## 🛠️ Implementasi Layanan

### 1️⃣ Validasi Email

**Nama Fungsi**: `ValidateEmail(email string)`

**Fungsi**: Memeriksa apakah format email yang dimasukkan pengguna valid

**Cara Kerja**:
- Menggunakan regex untuk validasi format email
- Pattern: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

**Error yang Mungkin Muncul**:
- ❌ "Format email tidak valid" (kode: 400)

### 2️⃣ Pembuat Token Aktivasi

**Nama Fungsi**: `GenerateActivationToken()`

**Fungsi**: Membuat token 6 digit untuk verifikasi

**Cara Kerja**:
- Menghasilkan token 6 digit berdasarkan timestamp
- Mengembalikan token dalam bentuk string

### 3️⃣ Validasi Data Registrasi

**Nama Fungsi**: `ValidateRegistrationInput(email, username, password, passwordConfirmation string)`

**Fungsi**: Memeriksa data registrasi pengguna baru

**Cara Kerja**:
- Memvalidasi format email
- Memastikan username tidak melebihi 30 karakter
- Memeriksa kecocokan password dengan konfirmasi

**Error yang Mungkin Muncul**:
- ❌ "Format email tidak valid" (kode: 400)
- ❌ "Username maksimal 30 karakter" (kode: 400)
- ❌ "Password tidak cocok dengan konfirmasi" (kode: 400)

### 4️⃣ Pengaman Password

**Nama Fungsi**: `HashPassword(password string)`

**Fungsi**: Mengubah password biasa menjadi bentuk terenkripsi

**Cara Kerja**:
- Menggunakan bcrypt dengan DefaultCost
- Menghasilkan hash password yang aman

**Error yang Mungkin Muncul**:
- ❌ "Gagal mengenkripsi password" (kode: 500)

### 4️⃣ HashString

**Nama Fungsi**: `HashString(input string)`

**Fungsi**: Mengubah string (password atau token) menjadi bentuk terenkripsi

**Cara Kerja**:
- Menggunakan bcrypt dengan DefaultCost
- Dapat digunakan untuk password maupun token
- Menghasilkan hash yang aman dan konsisten

**Error yang Mungkin Muncul**:
- ❌ "Gagal mengenkripsi string" (kode: 500)

### 5️⃣ VerifyHash

**Nama Fungsi**: `VerifyHash(hashedString, input string)`

**Fungsi**: Memeriksa apakah string yang dimasukkan cocok dengan hash yang tersimpan

**Cara Kerja**:
- Dapat memverifikasi hash dari password atau token
- Menggunakan bcrypt untuk perbandingan aman
- Mengembalikan error jika tidak cocok

**Error yang Mungkin Muncul**:
- ❌ "Hash tidak cocok" (kode: 401)

### 5️⃣ Pemeriksa Password

**Nama Fungsi**: `VerifyPassword(hashedPassword, password string)`

**Fungsi**: Memeriksa apakah password yang dimasukkan cocok dengan hash yang tersimpan

**Error yang Mungkin Muncul**:
- ❌ "Email atau password salah" (kode: 401)

### 6️⃣ Validasi Input Password

**Nama Fungsi**: `ValidatePasswordInput(password, passwordConfirmation string)`

**Fungsi**: Memvalidasi input password baru dengan konfirmasinya

**Error yang Mungkin Muncul**:
- ❌ "Password baru tidak cocok dengan konfirmasi" (kode: 400)

### 7️⃣ Pembuat Token

**Nama Fungsi**: `GenerateTokenPair(userID, email, userType string)`

**Fungsi**: Membuat pasangan access token dan refresh token

**Cara Kerja**:
- Menggunakan jwtService untuk membuat token
- Menghasilkan access token dan refresh token

**Error yang Mungkin Muncul**:
- ❌ "Gagal membuat access token" (kode: 500)
- ❌ "Gagal membuat refresh token" (kode: 500)

### 8️⃣ Validasi Token Claims

**Nama Fungsi**: `ValidateTokenClaims(token string)`

**Fungsi**: Memvalidasi dan mengekstrak klaim dari token JWT

**Error yang Mungkin Muncul**:
- ❌ "Token sudah kadaluarsa" (kode: 401)
- ❌ "Token tidak valid" (kode: 401)
- ❌ "Gagal memvalidasi token" (kode: 500)

### 9️⃣ Generate Access Token

**Nama Fungsi**: `GenerateAccessToken(userID, email, userType string)`

**Fungsi**: Membuat access token baru

**Error yang Mungkin Muncul**:
- ❌ "Gagal membuat access token baru" (kode: 500)

## ⚠️ Kode Error Umum

| Kode | Nama Error | Contoh Pesan Error |
|------|------------|-------------------|
| **400** | Bad Request | Format email tidak valid |
| **401** | Unauthorized | Token sudah kadaluarsa |
| **500** | Server Error | Gagal mengenkripsi password |

## 📘 Istilah Teknis

- **JWT (JSON Web Token)**: Format token yang digunakan untuk autentikasi
- **Access Token**: Token jangka pendek untuk mengakses API
- **Refresh Token**: Token jangka panjang untuk memperbaharui access token
- **Bcrypt**: Algoritma hash untuk mengamankan password
- **Claims**: Data yang disimpan dalam token JWT