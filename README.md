CMS + Portal

Sistem CMS dan Portal Publik berbasis Go (Gin) dan React (Vite + TypeScript) dengan arsitektur terpisah antara backend CMS dan frontend portal.

🏗 Arsitektur
Go (Gin Backend) React (Vite + TS)
├── CMS Panel (Server Render) ─┐ Portal Publik
│ - Authentication │ ├── Beranda
│ - Content Management │ ├── Berita
│ │ ├── Karier
└── REST API ─┘→ ├── Tentang
JSON API untuk frontend └── Kontak
🚀 Cara Menjalankan

1. Database (PostgreSQL via Docker)
   docker compose up -d
2. Backend (Go - Gin)
   cd backend
   cp .env.example .env
   go mod tidy
   go run ./cmd/app/

Backend akan berjalan di:

http://localhost:8080

CMS Panel:

http://localhost:8080/cms/login 3. Frontend (React - Vite + TypeScript)
cd frontend
cp .env.example .env
npm install
npm run dev

Akses portal:

http://localhost:5173
🔐 Akun Default (Development Only)

Digunakan hanya untuk kebutuhan development lokal

Role Email Password
Admin admin@example.local (set via seed / env)
🌐 Endpoint Utama
URL Keterangan
/cms/login Login CMS Panel
/api/v1/\* REST API Backend
http://localhost:5173 Portal Publik
📁 Struktur Proyek
cms-project/
├── backend/
│ ├── cmd/app/ # Entry point aplikasi
│ ├── config/ # Konfigurasi & database
│ ├── internal/
│ │ ├── domain/ # Model/entity
│ │ ├── repo/ # Query database
│ │ ├── handler/ # HTTP handler
│ │ └── middleware/ # Auth & middleware
│ ├── templates/ # HTML CMS (server-rendered)
│ │ ├── layouts/
│ │ ├── auth/
│ │ └── cms/
│ ├── static/ # Asset CSS, upload file
│ └── migrations/ # Database migration & seed
│
├── frontend/
│ └── src/
│ ├── api/ # API client
│ ├── components/ # UI components
│ └── pages/ # Halaman portal
│
└── docker-compose.yml # PostgreSQL setup
👥 Role System
Role Akses
super_admin Semua fitur
content_editor Kelola berita, banner, halaman
reviewer Review & approve konten
hr_recruitment Kelola lowongan & pelamar
🔄 Workflow Berita
Content Editor
↓ Submit
Reviewer / Admin
↓ Review & Approve
Published → Tampil di Portal Publik
⚙️ Tech Stack

Backend

Go (Gin Framework)
PostgreSQL
HTML Templates (CMS)
REST API

Frontend

React
Vite
TypeScript
Axios

DevOps

Docker Compose
📝 Catatan
Project ini menggunakan environment lokal untuk development
Semua konfigurasi sensitif disimpan di file .env
Struktur dibuat modular agar mudah dikembangkan
📌 Tujuan Project

Project ini dibuat untuk sistem informasi yayasan yang memiliki:

CMS internal untuk manajemen konten
Portal publik untuk user umum
Role-based access control (RBAC)
REST API untuk integrasi frontend
