package handler

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/yayasan/cms/internal/domain"
	"github.com/yayasan/cms/internal/repo"
)

type CMSNewsHandler struct {
	newsRepo *repo.NewsRepo
}

func NewCMSNewsHandler(nr *repo.NewsRepo) *CMSNewsHandler { return &CMSNewsHandler{nr} }

func (h *CMSNewsHandler) Index(c *gin.Context) {
	status := c.Query("status")
	search := c.Query("q")
	page := queryInt(c, "page", 1)
	news, total, _ := h.newsRepo.FindAll(domain.ListOptions{
		Page: page, Limit: 15, Status: status, Search: search,
	})
	c.HTML(http.StatusOK, "cms/news_list.html", withUserCtx(c, gin.H{
		"title": "Manajemen Berita", "active_menu": "news",
		"news": news, "total": total, "page": page,
		"status": status, "search": search, "flash": c.Query("flash"),
	}))
}

func (h *CMSNewsHandler) CreatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "cms/news_form.html", withUserCtx(c, gin.H{
		"title": "Tulis Berita", "active_menu": "news", "is_edit": false,
	}))
}

func (h *CMSNewsHandler) Create(c *gin.Context) {
	authorID := sessionUserID(c)
	title := c.PostForm("title")

	thumbnail := ""
	file, err := c.FormFile("thumbnail")
	if err == nil {
		ext := filepath.Ext(file.Filename)
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
		dst := "static/uploads/" + filename
		if c.SaveUploadedFile(file, dst) == nil {
			thumbnail = "/static/uploads/" + filename
		}
	}

	slug := generateSlug(title)
	if n, _ := h.newsRepo.FindBySlug(slug); n != nil {
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	metaTitle := c.PostForm("meta_title")
	if metaTitle == "" {
		metaTitle = title
	}

	h.newsRepo.Create(&domain.News{
		Title: title, Slug: slug,
		Content: c.PostForm("content"), Excerpt: c.PostForm("excerpt"),
		Thumbnail: thumbnail, Status: domain.StatusDraft,
		AuthorID: authorID, MetaTitle: metaTitle,
		MetaDescription: c.PostForm("meta_description"),
	})
	c.Redirect(http.StatusFound, "/cms/news?flash=created")
}

func (h *CMSNewsHandler) EditPage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	news, err := h.newsRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/news")
		return
	}
	c.HTML(http.StatusOK, "cms/news_form.html", withUserCtx(c, gin.H{
		"title": "Edit Berita", "active_menu": "news",
		"is_edit": true, "news": news,
	}))
}

func (h *CMSNewsHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	news, err := h.newsRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/news")
		return
	}

	news.Title = c.PostForm("title")
	news.Content = c.PostForm("content")
	news.Excerpt = c.PostForm("excerpt")
	news.MetaTitle = c.PostForm("meta_title")
	news.MetaDescription = c.PostForm("meta_description")

	file, err := c.FormFile("thumbnail")
	if err == nil {
		oldThumb := news.Thumbnail // simpan sebelum diganti
		ext := filepath.Ext(file.Filename)
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
		dst := "static/uploads/" + filename
		if c.SaveUploadedFile(file, dst) == nil {
			news.Thumbnail = "/static/uploads/" + filename
			removeFile(oldThumb) // hapus thumbnail lama
		}
	}

	h.newsRepo.Update(news)
	c.Redirect(http.StatusFound, "/cms/news?flash=updated")
}

func (h *CMSNewsHandler) SubmitReview(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	news, err := h.newsRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/news")
		return
	}
	news.Status = domain.StatusPendingReview
	h.newsRepo.Update(news)
	c.Redirect(http.StatusFound, "/cms/news?flash=submitted")
}

func (h *CMSNewsHandler) Approve(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	news, err := h.newsRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/news")
		return
	}
	now := time.Now()
	news.Status = domain.StatusPublished
	news.PublishedAt = &now
	h.newsRepo.Update(news)
	c.Redirect(http.StatusFound, "/cms/news?flash=approved")
}

func (h *CMSNewsHandler) Reject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	news, err := h.newsRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/news")
		return
	}
	news.Status = domain.StatusDraft
	h.newsRepo.Update(news)
	c.Redirect(http.StatusFound, "/cms/news?flash=rejected")
}

func (h *CMSNewsHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	news, err := h.newsRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/news?flash=deleted")
		return
	}
	thumbnail := news.Thumbnail // simpan sebelum hapus dari DB
	h.newsRepo.Delete(uint(id))
	removeFile(thumbnail) // hapus thumbnail dari disk
	c.Redirect(http.StatusFound, "/cms/news?flash=deleted")
}

// ── helpers ───────────────────────────────────────────────────────────────────

func generateSlug(title string) string {
	s := strings.ToLower(title)
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	clean := strings.Trim(b.String(), "-")
	for strings.Contains(clean, "--") {
		clean = strings.ReplaceAll(clean, "--", "-")
	}
	return clean
}

func queryInt(c *gin.Context, key string, def int) int {
	if v := c.Query(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
