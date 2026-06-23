package handler

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yayasan/cms/internal/domain"
	"github.com/yayasan/cms/internal/repo"
)

// APIHandler menyediakan data JSON untuk portal React
type APIHandler struct {
	newsRepo            *repo.NewsRepo
	volunteerRepo       *repo.VolunteerRepo
	applicantRepo       *repo.ApplicantRepo
	pageRepo            *repo.PageRepo
	bannerRepo          *repo.BannerRepo
	galleryRepo         *repo.GalleryRepo
	galleryCategoryRepo *repo.GalleryCategoryRepo
	settingRepo         *repo.SettingRepo
}

func NewAPIHandler(
	nr *repo.NewsRepo,
	vr *repo.VolunteerRepo,
	ar *repo.ApplicantRepo,
	pr *repo.PageRepo,
	br *repo.BannerRepo,
	gr *repo.GalleryRepo,
	gcr *repo.GalleryCategoryRepo,
	sr *repo.SettingRepo,
) *APIHandler {
	return &APIHandler{nr, vr, ar, pr, br, gr, gcr, sr}
}

// ── Berita ────────────────────────────────────────────────────────────────────

// GET /api/v1/news
func (h *APIHandler) NewsList(c *gin.Context) {
	opts := domain.ListOptions{
		Status: "published",
		Page:   queryInt(c, "page", 1),
		Limit:  queryInt(c, "limit", 9),
		Search: c.Query("q"),
	}
	news, total, _ := h.newsRepo.FindAll(opts)
	c.JSON(http.StatusOK, gin.H{
		"data": news,
		"meta": gin.H{"total": total, "page": opts.Page, "limit": opts.Limit},
	})
}

// GET /api/v1/news/:slug
func (h *APIHandler) NewsDetail(c *gin.Context) {
	news, err := h.newsRepo.FindBySlug(c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}
	h.newsRepo.IncrementView(news.ID)
	c.JSON(http.StatusOK, gin.H{"data": news})
}

// ── Relawan ───────────────────────────────────────────────────────────────────

// GET /api/v1/volunteers
func (h *APIHandler) VolunteerList(c *gin.Context) {
	opts := domain.ListOptions{
		Status: "active",
		Page:   queryInt(c, "page", 1),
		Limit:  queryInt(c, "limit", 10),
	}
	volunteers, total, _ := h.volunteerRepo.FindAll(opts)
	c.JSON(http.StatusOK, gin.H{
		"data": volunteers,
		"meta": gin.H{"total": total, "page": opts.Page, "limit": opts.Limit},
	})
}

// GET /api/v1/volunteers/:slug
func (h *APIHandler) VolunteerDetail(c *gin.Context) {
	volunteer, err := h.volunteerRepo.FindBySlug(c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rekrutmen tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": volunteer})
}

// POST /api/v1/volunteers/:slug/apply — form pendaftaran relawan dari portal
func (h *APIHandler) Apply(c *gin.Context) {
	volunteer, err := h.volunteerRepo.FindBySlug(c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rekrutmen tidak ditemukan"})
		return
	}

	applicant := &domain.Applicant{
		VolunteerID: volunteer.ID,
		FullName:    c.PostForm("full_name"),
		Email:       c.PostForm("email"),
		Phone:       c.PostForm("phone"),
		City:        c.PostForm("city"),
		Occupation:  c.PostForm("occupation"),
		Motivation:  c.PostForm("motivation"),
		Status:      domain.AppApplied,
	}

	// Upload CV
	file, err := c.FormFile("cv")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_cv" + ext
		dst := "static/uploads/cv/" + filename
		if c.SaveUploadedFile(file, dst) == nil {
			applicant.CvPath = dst
		}
	}

	if err := h.applicantRepo.Create(applicant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pendaftaran"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pendaftaran berhasil dikirim, kami akan menghubungi Anda",
	})
}

// ── Galeri ────────────────────────────────────────────────────────────────────

// GET /api/v1/gallery
func (h *APIHandler) GalleryList(c *gin.Context) {
	opts := domain.ListOptions{
		Page:     queryInt(c, "page", 1),
		Limit:    queryInt(c, "limit", 12),
		Category: c.Query("category"), // slug kategori
	}
	items, total, _ := h.galleryRepo.FindActive(opts)
	c.JSON(http.StatusOK, gin.H{
		"data": items,
		"meta": gin.H{"total": total, "page": opts.Page, "limit": opts.Limit},
	})
}

// GET /api/v1/gallery/categories
func (h *APIHandler) GalleryCategories(c *gin.Context) {
	categories, _ := h.galleryCategoryRepo.FindAll()
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// ── Halaman Statis ────────────────────────────────────────────────────────────

// GET /api/v1/pages/:slug
func (h *APIHandler) PageDetail(c *gin.Context) {
	page, err := h.pageRepo.FindBySlug(c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Halaman tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": page})
}

// ── Banner ────────────────────────────────────────────────────────────────────

// GET /api/v1/banners
func (h *APIHandler) Banners(c *gin.Context) {
	banners, _ := h.bannerRepo.FindActive()
	c.JSON(http.StatusOK, gin.H{"data": banners})
}

// ── Site Settings (konten dinamis: hero, statistik, footer, kontak, dll) ──────

// GET /api/v1/settings
func (h *APIHandler) Settings(c *gin.Context) {
	m, err := h.settingRepo.AsMap()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat pengaturan situs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

// ── Home ──────────────────────────────────────────────────────────────────────

// GET /api/v1/home
func (h *APIHandler) Home(c *gin.Context) {
	banners, _ := h.bannerRepo.FindActive()
	latestNews, _ := h.newsRepo.FindPublished(6)
	latestVolunteers, _ := h.volunteerRepo.FindActive(4)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"banners":           banners,
			"latest_news":       latestNews,
			"latest_volunteers": latestVolunteers,
		},
	})
}
