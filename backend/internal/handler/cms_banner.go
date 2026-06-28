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

type CMSBannerHandler struct {
	bannerRepo *repo.BannerRepo
}

func NewCMSBannerHandler(br *repo.BannerRepo) *CMSBannerHandler {
	return &CMSBannerHandler{br}
}

// GET /cms/banners
func (h *CMSBannerHandler) Index(c *gin.Context) {
	banners, _ := h.bannerRepo.FindAll()
	c.HTML(http.StatusOK, "cms/banner_list.html", withUserCtx(c, gin.H{
		"title": "Manajemen Banner", "active_menu": "banners",
		"banners": banners, "flash": c.Query("flash"),
	}))
}

// GET /cms/banners/create
func (h *CMSBannerHandler) CreatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "cms/banner_form.html", withUserCtx(c, gin.H{
		"title": "Tambah Banner", "active_menu": "banners", "is_edit": false,
	}))
}

// POST /cms/banners/create
func (h *CMSBannerHandler) Create(c *gin.Context) {
	order, _ := strconv.Atoi(c.PostForm("order_num"))

	imagePath := ""
	file, err := c.FormFile("image")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
		dst := "static/uploads/banners/" + filename
		if c.SaveUploadedFile(file, dst) == nil {
			imagePath = "/static/uploads/banners/" + filename
		}
	}

	if imagePath == "" {
		c.HTML(http.StatusOK, "cms/banner_form.html", withUserCtx(c, gin.H{
			"title": "Tambah Banner", "active_menu": "banners",
			"is_edit": false, "error": "Gambar banner wajib diupload",
		}))
		return
	}

	h.bannerRepo.Create(&domain.Banner{
		Title: c.PostForm("title"), ImagePath: imagePath,
		LinkURL: c.PostForm("link_url"), OrderNum: order,
		IsActive: c.PostForm("is_active") == "1",
	})
	c.Redirect(http.StatusFound, "/cms/banners?flash=created")
}

// GET /cms/banners/:id/edit
func (h *CMSBannerHandler) EditPage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	banner, err := h.bannerRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/banners")
		return
	}
	c.HTML(http.StatusOK, "cms/banner_form.html", withUserCtx(c, gin.H{
		"title": "Edit Banner", "active_menu": "banners",
		"is_edit": true, "banner": banner,
	}))
}

// POST /cms/banners/:id/edit
func (h *CMSBannerHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	banner, err := h.bannerRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/banners")
		return
	}

	order, _ := strconv.Atoi(c.PostForm("order_num"))
	banner.Title = c.PostForm("title")
	banner.LinkURL = c.PostForm("link_url")
	banner.OrderNum = order
	banner.IsActive = c.PostForm("is_active") == "1"

	file, err := c.FormFile("image")
	if err == nil {
		oldPath := banner.ImagePath // simpan path lama sebelum diganti
		ext := filepath.Ext(file.Filename)
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
		dst := "static/uploads/banners/" + filename
		if c.SaveUploadedFile(file, dst) == nil {
			banner.ImagePath = "/static/uploads/banners/" + filename
			removeFile(oldPath) // hapus gambar lama dari disk
		}
	}

	h.bannerRepo.Update(banner)
	c.Redirect(http.StatusFound, "/cms/banners?flash=updated")
}

// POST /cms/banners/:id/toggle
func (h *CMSBannerHandler) Toggle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	banner, err := h.bannerRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/banners")
		return
	}
	banner.IsActive = !banner.IsActive
	h.bannerRepo.Update(banner)
	c.Redirect(http.StatusFound, "/cms/banners?flash=toggled")
}

// POST /cms/banners/:id/delete
func (h *CMSBannerHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	banner, err := h.bannerRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/banners?flash=deleted")
		return
	}
	imagePath := banner.ImagePath // simpan sebelum hapus dari DB
	h.bannerRepo.Delete(uint(id))
	removeFile(imagePath) // hapus file dari disk
	c.Redirect(http.StatusFound, "/cms/banners?flash=deleted")
}
