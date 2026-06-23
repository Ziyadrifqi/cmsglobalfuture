package domain

import (
	"time"

	"gorm.io/gorm"
)

// ── Role ─────────────────────────────────────────────────────────────────────

type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null"`
	DisplayName string `gorm:"not null"`
	CreatedAt   time.Time
}

const (
	RoleSuperAdmin    = "super_admin"
	RoleContentEditor = "content_editor"
	RoleHR            = "hr_recruitment"
	RoleReviewer      = "reviewer"
)

// ── User ─────────────────────────────────────────────────────────────────────

type User struct {
	gorm.Model
	Name               string `gorm:"not null"`
	Email              string `gorm:"uniqueIndex;not null"`
	Password           string `gorm:"not null"`
	Avatar             string
	MustChangePassword bool `gorm:"default:false"`
	IsActive           bool `gorm:"default:true"`
	RoleID             uint `gorm:"not null"`
	Role               Role `gorm:"foreignKey:RoleID"`
}

// ── News ─────────────────────────────────────────────────────────────────────

type NewsStatus string

const (
	StatusDraft         NewsStatus = "draft"
	StatusPendingReview NewsStatus = "pending_review"
	StatusPublished     NewsStatus = "published"
)

type NewsCategory struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Slug      string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
}

type News struct {
	gorm.Model
	Title           string `gorm:"not null"`
	Slug            string `gorm:"uniqueIndex;not null"`
	Content         string `gorm:"type:text"`
	Excerpt         string
	Thumbnail       string
	Status          NewsStatus `gorm:"default:'draft'"`
	PublishedAt     *time.Time
	CategoryID      *uint
	Category        *NewsCategory `gorm:"foreignKey:CategoryID"`
	AuthorID        uint
	Author          User `gorm:"foreignKey:AuthorID"`
	MetaTitle       string
	MetaDescription string
	ViewCount       int `gorm:"default:0"`
}

// ── Volunteer (sebelumnya: Job) ───────────────────────────────────────────────

type VolunteerStatus string

const (
	VolunteerDraft  VolunteerStatus = "draft"
	VolunteerActive VolunteerStatus = "active"
	VolunteerClosed VolunteerStatus = "closed"
)

type Volunteer struct {
	gorm.Model
	Title           string `gorm:"not null"`
	Slug            string `gorm:"uniqueIndex;not null"`
	Division        string // sebelumnya: Department
	Location        string
	Type            string          `gorm:"default:'regular'"` // regular, event, remote, training
	Description     string          `gorm:"type:text"`
	Requirements    string          `gorm:"type:text"`
	Benefits        string          `gorm:"type:text"`
	Status          VolunteerStatus `gorm:"default:'draft'"`
	ClosedAt        *time.Time
	CreatedByID     uint
	CreatedBy       User `gorm:"foreignKey:CreatedByID"`
	MetaTitle       string
	MetaDescription string
}

// ── Applicant (Pendaftar Relawan) ─────────────────────────────────────────────

type ApplicantStatus string

const (
	AppApplied   ApplicantStatus = "applied"
	AppScreening ApplicantStatus = "screening"
	AppInterview ApplicantStatus = "interview"
	AppAccepted  ApplicantStatus = "accepted"
	AppRejected  ApplicantStatus = "rejected"
)

type Applicant struct {
	gorm.Model
	VolunteerID uint      `gorm:"not null"`
	Volunteer   Volunteer `gorm:"foreignKey:VolunteerID"`
	FullName    string    `gorm:"not null"`
	Email       string    `gorm:"not null"`
	Phone       string
	City        string
	Occupation  string
	CvPath      string
	Motivation  string          `gorm:"type:text"`
	Status      ApplicantStatus `gorm:"default:'applied'"`
	Notes       string          `gorm:"type:text"`
}

// ── Banner ────────────────────────────────────────────────────────────────────

type Banner struct {
	gorm.Model
	Title     string `gorm:"not null"`
	ImagePath string `gorm:"not null"`
	LinkURL   string
	OrderNum  int  `gorm:"default:0"`
	IsActive  bool `gorm:"default:true"`
}

// ── Static Page ───────────────────────────────────────────────────────────────

type Page struct {
	gorm.Model
	Slug            string `gorm:"uniqueIndex;not null"`
	Title           string `gorm:"not null"`
	Content         string `gorm:"type:text"`
	MetaTitle       string
	MetaDescription string
}

// ── Gallery ───────────────────────────────────────────────────────────────────

type GalleryType string

const (
	GalleryImage GalleryType = "image"
	GalleryVideo GalleryType = "video"
)

// Kategori bawaan galeri — bisa ditambah via seed atau CMS nanti.
const (
	GalleryCatMangrove  = "Penanaman Mangrove"
	GalleryCatBeach     = "Bersih Pantai"
	GalleryCatEducation = "Edukasi"
	GalleryCatWorkshop  = "Workshop"
	GalleryCatEvent     = "Event Nasional"
)

type GalleryCategory struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	Slug      string `gorm:"uniqueIndex;not null"`
	OrderNum  int    `gorm:"default:0"`
	CreatedAt time.Time
}

type GalleryItem struct {
	gorm.Model
	Title        string      `gorm:"not null"`
	Type         GalleryType `gorm:"not null;default:'image'"`
	URL          string      `gorm:"not null"` // URL file asli (gambar/video)
	ThumbnailURL string      // thumbnail untuk video
	Caption      string
	CategoryID   *uint
	Category     *GalleryCategory `gorm:"foreignKey:CategoryID"`
	IsActive     bool             `gorm:"default:true"`
	OrderNum     int              `gorm:"default:0"`
	UploadedByID uint
	UploadedBy   User `gorm:"foreignKey:UploadedByID"`
}

// ── ListOptions (untuk filter + pagination di repo) ───────────────────────────

type ListOptions struct {
	Page     int
	Limit    int
	Status   string
	Search   string
	Category string
}
