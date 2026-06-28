package handler

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yayasan/cms/internal/domain"
	"github.com/yayasan/cms/internal/repo"
)

type CMSGalleryHandler struct {
	galleryRepo  *repo.GalleryRepo
	categoryRepo *repo.GalleryCategoryRepo
}

func NewCMSGalleryHandler(gr *repo.GalleryRepo, cr *repo.GalleryCategoryRepo) *CMSGalleryHandler {
	return &CMSGalleryHandler{gr, cr}
}

func (h *CMSGalleryHandler) Index(c *gin.Context) {
	items, total, _ := h.galleryRepo.FindAll(domain.ListOptions{
		Page: queryInt(c, "page", 1), Limit: 20, Category: c.Query("category"),
	})
	categories, _ := h.categoryRepo.FindAll()
	c.HTML(http.StatusOK, "cms/gallery_list.html", withUserCtx(c, gin.H{
		"title": "Manajemen Galeri", "active_menu": "gallery",
		"items": items, "total": total, "categories": categories,
		"active_category": c.Query("category"), "flash": c.Query("flash"),
	}))
}

func (h *CMSGalleryHandler) CreatePage(c *gin.Context) {
	categories, _ := h.categoryRepo.FindAll()
	c.HTML(http.StatusOK, "cms/gallery_form.html", withUserCtx(c, gin.H{
		"title": "Tambah Media Galeri", "active_menu": "gallery",
		"categories": categories, "is_edit": false,
	}))
}

func (h *CMSGalleryHandler) Create(c *gin.Context) {
	uploadedBy := sessionUserID(c)
	mediaType := domain.GalleryType(c.PostForm("type"))
	if mediaType != domain.GalleryImage && mediaType != domain.GalleryVideo {
		mediaType = domain.GalleryImage
	}
	categoryID := parseUintPtr(c.PostForm("category_id"))
	order, _ := strconv.Atoi(c.PostForm("order_num"))

	mediaURL, thumbURL, errMsg := h.handleUpload(c, mediaType)
	if errMsg != "" {
		categories, _ := h.categoryRepo.FindAll()
		c.HTML(http.StatusOK, "cms/gallery_form.html", withUserCtx(c, gin.H{
			"title": "Tambah Media Galeri", "active_menu": "gallery",
			"categories": categories, "is_edit": false, "error": errMsg,
		}))
		return
	}

	h.galleryRepo.Create(&domain.GalleryItem{
		Title: c.PostForm("title"), Type: mediaType,
		URL: mediaURL, ThumbnailURL: thumbURL,
		Caption: c.PostForm("caption"), CategoryID: categoryID,
		IsActive: c.PostForm("is_active") != "0",
		OrderNum: order, UploadedByID: uploadedBy,
	})
	c.Redirect(http.StatusFound, "/cms/gallery?flash=created")
}

func (h *CMSGalleryHandler) EditPage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, err := h.galleryRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/gallery")
		return
	}
	categories, _ := h.categoryRepo.FindAll()
	c.HTML(http.StatusOK, "cms/gallery_form.html", withUserCtx(c, gin.H{
		"title": "Edit Media Galeri", "active_menu": "gallery",
		"categories": categories, "is_edit": true, "item": item,
	}))
}

func (h *CMSGalleryHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, err := h.galleryRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/gallery")
		return
	}

	item.Title = c.PostForm("title")
	item.Caption = c.PostForm("caption")
	item.CategoryID = parseUintPtr(c.PostForm("category_id"))
	item.IsActive = c.PostForm("is_active") != "0"
	order, _ := strconv.Atoi(c.PostForm("order_num"))
	item.OrderNum = order

	newURL, newThumb, _ := h.handleUpload(c, item.Type)

	if newURL != "" {
		removeFile(item.URL) // hapus file media lama
		item.URL = newURL
	}
	if newThumb != "" {
		removeFile(item.ThumbnailURL) // hapus thumbnail lama
		item.ThumbnailURL = newThumb
	}

	h.galleryRepo.Update(item)
	c.Redirect(http.StatusFound, "/cms/gallery?flash=updated")
}

func (h *CMSGalleryHandler) Toggle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, err := h.galleryRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/gallery")
		return
	}
	item.IsActive = !item.IsActive
	h.galleryRepo.Update(item)
	c.Redirect(http.StatusFound, "/cms/gallery?flash=toggled")
}

func (h *CMSGalleryHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, err := h.galleryRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/gallery?flash=deleted")
		return
	}
	// Simpan path sebelum hapus dari DB
	mediaURL := item.URL
	thumbURL := item.ThumbnailURL

	h.galleryRepo.Delete(uint(id))

	// Hapus file dari disk setelah data terhapus dari DB
	removeFile(mediaURL)
	removeFile(thumbURL)

	c.Redirect(http.StatusFound, "/cms/gallery?flash=deleted")
}

// ── Category CRUD ─────────────────────────────────────────────────────────────

func (h *CMSGalleryHandler) CategoryIndex(c *gin.Context) {
	categories, _ := h.categoryRepo.FindAll()
	c.HTML(http.StatusOK, "cms/gallery_category_list.html", withUserCtx(c, gin.H{
		"title": "Kategori Galeri", "active_menu": "gallery",
		"categories": categories, "flash": c.Query("flash"),
	}))
}

func (h *CMSGalleryHandler) CategoryCreate(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	if name == "" {
		c.Redirect(http.StatusFound, "/cms/gallery/categories?flash=error_empty")
		return
	}
	order, _ := strconv.Atoi(c.PostForm("order_num"))
	h.categoryRepo.Create(&domain.GalleryCategory{
		Name: name, Slug: generateSlug(name), OrderNum: order,
	})
	c.Redirect(http.StatusFound, "/cms/gallery/categories?flash=created")
}

func (h *CMSGalleryHandler) CategoryUpdate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	cat, err := h.categoryRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/gallery/categories")
		return
	}
	cat.Name = strings.TrimSpace(c.PostForm("name"))
	order, _ := strconv.Atoi(c.PostForm("order_num"))
	cat.OrderNum = order
	h.categoryRepo.Update(cat)
	c.Redirect(http.StatusFound, "/cms/gallery/categories?flash=updated")
}

func (h *CMSGalleryHandler) CategoryDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.categoryRepo.Delete(uint(id))
	c.Redirect(http.StatusFound, "/cms/gallery/categories?flash=deleted")
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (h *CMSGalleryHandler) handleUpload(c *gin.Context, mediaType domain.GalleryType) (mediaURL, thumbURL, errMsg string) {
	file, err := c.FormFile("media_file")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
		dst := "static/uploads/gallery/" + filename
		if err2 := c.SaveUploadedFile(file, dst); err2 == nil {
			mediaURL = "/static/uploads/gallery/" + filename
		} else {
			errMsg = "Gagal menyimpan file media"
			return
		}
	}
	thumb, err := c.FormFile("thumbnail_file")
	if err == nil {
		ext := filepath.Ext(thumb.Filename)
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_thumb" + ext
		dst := "static/uploads/gallery/" + filename
		if err2 := c.SaveUploadedFile(thumb, dst); err2 == nil {
			thumbURL = "/static/uploads/gallery/" + filename
		}
	}
	return
}

func parseUintPtr(s string) *uint {
	if s == "" || s == "0" {
		return nil
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil
	}
	v := uint(n)
	return &v
}
