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

type CMSVolunteerHandler struct {
	volunteerRepo *repo.VolunteerRepo
	applicantRepo *repo.ApplicantRepo
}

func NewCMSVolunteerHandler(vr *repo.VolunteerRepo, ar *repo.ApplicantRepo) *CMSVolunteerHandler {
	return &CMSVolunteerHandler{vr, ar}
}

// GET /cms/volunteers
func (h *CMSVolunteerHandler) Index(c *gin.Context) {
	volunteers, total, _ := h.volunteerRepo.FindAll(domain.ListOptions{
		Page:   queryInt(c, "page", 1),
		Limit:  15,
		Status: c.Query("status"),
		Search: c.Query("q"),
	})
	c.HTML(http.StatusOK, "cms/volunteer_list.html", withUserCtx(c, gin.H{
		"title":       "Manajemen Relawan",
		"active_menu": "volunteers",
		"volunteers":  volunteers,
		"total":       total,
		"flash":       c.Query("flash"),
	}))
}

// GET /cms/volunteers/create
func (h *CMSVolunteerHandler) CreatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "cms/volunteer_form.html", withUserCtx(c, gin.H{
		"title":       "Buka Rekrutmen Relawan",
		"active_menu": "volunteers",
		"is_edit":     false,
	}))
}

// POST /cms/volunteers/create
func (h *CMSVolunteerHandler) Create(c *gin.Context) {
	createdBy := sessionUserID(c)
	title := c.PostForm("title")
	slug := generateSlug(title)
	if v, _ := h.volunteerRepo.FindBySlugAny(slug); v != nil {
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	volunteer := &domain.Volunteer{
		Title:        title,
		Slug:         slug,
		Division:     c.PostForm("division"),
		Location:     c.PostForm("location"),
		Type:         c.PostForm("type"),
		Description:  c.PostForm("description"),
		Requirements: c.PostForm("requirements"),
		Benefits:     c.PostForm("benefits"),
		Status:       domain.VolunteerActive,
		CreatedByID:  createdBy,
	}
	h.volunteerRepo.Create(volunteer)
	c.Redirect(http.StatusFound, "/cms/volunteers?flash=created")
}

// GET /cms/volunteers/:id/edit
func (h *CMSVolunteerHandler) EditPage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	volunteer, err := h.volunteerRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/volunteers")
		return
	}
	c.HTML(http.StatusOK, "cms/volunteer_form.html", withUserCtx(c, gin.H{
		"title":       "Edit Rekrutmen",
		"active_menu": "volunteers",
		"is_edit":     true,
		"volunteer":   volunteer,
	}))
}

// POST /cms/volunteers/:id/edit
func (h *CMSVolunteerHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	volunteer, err := h.volunteerRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/volunteers")
		return
	}
	volunteer.Title = c.PostForm("title")
	volunteer.Division = c.PostForm("division")
	volunteer.Location = c.PostForm("location")
	volunteer.Type = c.PostForm("type")
	volunteer.Description = c.PostForm("description")
	volunteer.Requirements = c.PostForm("requirements")
	volunteer.Benefits = c.PostForm("benefits")
	h.volunteerRepo.Update(volunteer)
	c.Redirect(http.StatusFound, "/cms/volunteers?flash=updated")
}

// POST /cms/volunteers/:id/close
func (h *CMSVolunteerHandler) Close(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	volunteer, err := h.volunteerRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/volunteers")
		return
	}
	now := time.Now()
	volunteer.Status = domain.VolunteerClosed
	volunteer.ClosedAt = &now
	h.volunteerRepo.Update(volunteer)
	c.Redirect(http.StatusFound, "/cms/volunteers?flash=closed")
}

// POST /cms/volunteers/:id/delete
func (h *CMSVolunteerHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.volunteerRepo.Delete(uint(id))
	c.Redirect(http.StatusFound, "/cms/volunteers?flash=deleted")
}

// GET /cms/volunteers/:id/applicants
func (h *CMSVolunteerHandler) Applicants(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	volunteer, err := h.volunteerRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/volunteers")
		return
	}
	applicants, _ := h.applicantRepo.FindByVolunteer(uint(id))
	c.HTML(http.StatusOK, "cms/applicant_list.html", withUserCtx(c, gin.H{
		"title":       "Pendaftar — " + volunteer.Title,
		"active_menu": "volunteers",
		"volunteer":   volunteer,
		"applicants":  applicants,
		"flash":       c.Query("flash"),
	}))
}

// POST /cms/volunteers/applicants/:id/status
func (h *CMSVolunteerHandler) UpdateApplicantStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	applicant, err := h.applicantRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/volunteers")
		return
	}
	applicant.Status = domain.ApplicantStatus(c.PostForm("status"))
	applicant.Notes = c.PostForm("notes")
	h.applicantRepo.Update(applicant)
	c.Redirect(http.StatusFound,
		"/cms/volunteers/"+strconv.Itoa(int(applicant.VolunteerID))+"/applicants?flash=updated")
}

// GET /cms/volunteers/applicants/:id/cv
func (h *CMSVolunteerHandler) DownloadCV(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	applicant, err := h.applicantRepo.FindByID(uint(id))
	if err != nil || applicant.CvPath == "" {
		c.String(http.StatusNotFound, "CV tidak ditemukan")
		return
	}
	c.FileAttachment(applicant.CvPath, filepath.Base(applicant.CvPath))
}
