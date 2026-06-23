package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yayasan/cms/internal/repo"
)

type CMSPageHandler struct {
	pageRepo *repo.PageRepo
}

func NewCMSPageHandler(pr *repo.PageRepo) *CMSPageHandler {
	return &CMSPageHandler{pr}
}

// GET /cms/pages
func (h *CMSPageHandler) Index(c *gin.Context) {
	pages, _ := h.pageRepo.FindAll()
	c.HTML(http.StatusOK, "cms/page_list.html", withUserCtx(c, gin.H{
		"title":       "Halaman Statis",
		"active_menu": "pages",
		"pages":       pages,
		"flash":       c.Query("flash"),
	}))
}

// GET /cms/pages/:id/edit
func (h *CMSPageHandler) EditPage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	page, err := h.pageRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/pages")
		return
	}
	c.HTML(http.StatusOK, "cms/page_form.html", withUserCtx(c, gin.H{
		"title":       "Edit Halaman — " + page.Title,
		"active_menu": "pages",
		"page":        page,
	}))
}

// POST /cms/pages/:id/edit
func (h *CMSPageHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	page, err := h.pageRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/pages")
		return
	}
	page.Title = c.PostForm("title")
	page.Content = c.PostForm("content")
	page.MetaTitle = c.PostForm("meta_title")
	page.MetaDescription = c.PostForm("meta_description")
	h.pageRepo.Update(page)
	c.Redirect(http.StatusFound, "/cms/pages?flash=updated")
}
