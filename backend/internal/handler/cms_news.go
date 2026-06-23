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

func NewCMSNewsHandler(nr *repo.NewsRepo) *CMSNewsHandler {
	return &CMSNewsHandler{nr}
}

// GET /cms/news
func (h *CMSNewsHandler) Index(c *gin.Context) {
	status := c.Query("status")
	search := c.Query("q")
	page   := queryInt(c, "page", 1)

	news, total, _ := h.newsRepo.FindAll(domain.ListOptions{
		Page: page, Limit: 15, Status: status, Search: search,
	})

	// BUG FIX: withUserCtx() agar session_user_name/role terkirim ke template
	c.HTML(http.StatusOK, "cms/news_list.html", withUserCtx(c, gin.H{
		"title":       "Manajemen Berita",
		"active_menu": "news",
		"news":        news,
		"total":       total,
		"page":        page,
		"status":      status,
		"search":      search,
		"flash":       c.Query("flash"),
	}))
}

// GET /cms/news/create
func (h *CMSNewsHandler) CreatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "cms/news_form.html", withUserCtx(c, gin.H{
		"title":       "Tulis Berita",
		"active_menu": "news",
		"is_edit":     false,
	}))
}

// POST /cms/news/create
func (h *CMSNewsHandler) Create(c *gin.Context) {
	// BUG FIX #2: pakai sessionUserID() agar tidak panic
	authorID := sessionUserID(c)

	title      := c.PostForm("title")
	content    := c.PostForm("content")
	excerpt    := c.PostForm("excerpt")
	metaTitle  := c.PostForm("meta_title")
	metaDesc   := c.PostForm("meta_description")

	thumbnail := ""
	file, err := c.FormFile("thumbnail")
	if err == nil {
		ext      := filepath.Ext(file.Filename)
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
		dst      := "static/uploads/" + filename
		if c.SaveUploadedFile(file, dst) == nil {
			thumbnail = "/static/uploads/" + filename
		}
	}

	slug := generateSlug(title)
	if n, _ := h.newsRepo.FindBySlug(slug); n != nil {
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	news := &domain.News{
		Title:           title,
		Slug:            slug,
		Content:         content,
		Excerpt:         excerpt,
		Thumbnail:       thumbnail,
		Status:          domain.StatusDraft,
		AuthorID:        authorID,
		MetaTitle:       metaTitle,
		MetaDescription: metaDesc,
	}
	if news.MetaTitle == "" {
		news.MetaTitle = title
	}

	h.newsRepo.Create(news)
	c.Redirect(http.StatusFound, "/cms/news?flash=created")
}

// GET /cms/news/:id/edit
func (h *CMSNewsHandler) EditPage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	news, err := h.newsRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/news")
		return
	}
	c.HTML(http.StatusOK, "cms/news_form.html", withUserCtx(c, gin.H{
		"title":       "Edit Berita",
		"active_menu": "news",
		"is_edit":     true,
		"news":        news,
	}))
}

// POST /cms/news/:id/edit
func (h *CMSNewsHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	news, err := h.newsRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/news")
		return
	}

	news.Title           = c.PostForm("title")
	news.Content         = c.PostForm("content")
	news.Excerpt         = c.PostForm("excerpt")
	news.MetaTitle       = c.PostForm("meta_title")
	news.MetaDescription = c.PostForm("meta_description")

	file, err := c.FormFile("thumbnail")
	if err == nil {
		ext      := filepath.Ext(file.Filename)
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
		dst      := "static/uploads/" + filename
		if c.SaveUploadedFile(file, dst) == nil {
			news.Thumbnail = "/static/uploads/" + filename
		}
	}

	h.newsRepo.Update(news)
	c.Redirect(http.StatusFound, "/cms/news?flash=updated")
}

// POST /cms/news/:id/submit-review
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

// POST /cms/news/:id/approve
func (h *CMSNewsHandler) Approve(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	news, err := h.newsRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/news")
		return
	}
	now := time.Now()
	news.Status      = domain.StatusPublished
	news.PublishedAt = &now
	h.newsRepo.Update(news)
	c.Redirect(http.StatusFound, "/cms/news?flash=approved")
}

// POST /cms/news/:id/reject
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

// POST /cms/news/:id/delete
func (h *CMSNewsHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.newsRepo.Delete(uint(id))
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
