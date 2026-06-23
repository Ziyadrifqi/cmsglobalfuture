# Yayasan CMS + Portal

## Arsitektur

```
Go (Gin)                          React (Vite + TS)
├── CMS Panel (/cms/...)    ─┐    Portal Publik
│   Render HTML template    │    ├── Beranda
│   Login pakai session     │    ├── Berita
│                           │    ├── Karier + Form Lamaran
└── REST API (/api/v1/...)  ─┘─→  ├── Tentang
    Kirim JSON ke React          └── Kontak
```

---

## Cara Menjalankan

### 1. Database
```bash
docker compose up -d
```

### 2. Backend (Go)
```bash
cd backend
cp .env.example .env
go mod tidy
go run ./cmd/app/

# Output:
# ✓ Database terhubung
# ✓ Admin default: admin@yayasan.local / Admin@123
# Server jalan di http://localhost:8080
# CMS Panel: http://localhost:8080/cms/login
```

### 3. Frontend Portal (React)
```bash
cd frontend
cp .env.example .env
npm install
npm run dev

# Buka http://localhost:5173
```

---

## Akun Default

| Email                  | Password   | Role        |
|------------------------|------------|-------------|
| admin@yayasan.local    | Admin@123  | Super Admin |

---

## URL Penting

| URL                              | Keterangan               |
|----------------------------------|--------------------------|
| http://localhost:8080/cms/login  | Login CMS panel (Go)     |
| http://localhost:8080/api/v1/... | REST API                 |
| http://localhost:5173            | Portal publik (React)    |

---

## Struktur Folder

```
cms-project/
├── backend/                    ← Go (Gin)
│   ├── cmd/app/main.go         ← Entry point
│   ├── config/                 ← Konfigurasi & DB
│   ├── internal/
│   │   ├── domain/models.go    ← Semua struct entitas
│   │   ├── repo/               ← Query database
│   │   ├── handler/            ← HTTP handler
│   │   └── middleware/         ← Auth session & JWT
│   ├── templates/              ← HTML template CMS
│   │   ├── layouts/            ← Layout & sidebar
│   │   ├── auth/               ← Halaman login
│   │   └── cms/                ← Semua halaman CMS
│   ├── static/                 ← File statis (CSS, upload)
│   └── migrations/             ← Auto migrate + seed
│
├── frontend/                   ← React (Vite + TS)
│   └── src/
│       ├── api/index.ts        ← Semua panggilan API ke Go
│       ├── components/layout/  ← Navbar & Footer
│       └── pages/              ← Halaman portal
│
└── docker-compose.yml          ← PostgreSQL lokal
```

---

## 4 Role

| Role             | Akses CMS                           |
|------------------|-------------------------------------|
| super_admin      | Semua fitur                         |
| content_editor   | Berita, Banner, Halaman             |
| reviewer         | Review & approve berita             |
| hr_recruitment   | Lowongan & Pelamar                  |

## Approval Flow Berita

```
Content Editor
      ↓ Ajukan Review
Reviewer / Super Admin
      ↓ Approve
    Published (tampil di portal)
```
