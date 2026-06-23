package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yayasan/cms/internal/domain"
	"github.com/yayasan/cms/internal/repo"
)

type CMSDashboardHandler struct {
	newsRepo      *repo.NewsRepo
	volunteerRepo *repo.VolunteerRepo
	applicantRepo *repo.ApplicantRepo
	userRepo      *repo.UserRepo
}

func NewCMSDashboardHandler(
	nr *repo.NewsRepo,
	vr *repo.VolunteerRepo,
	ar *repo.ApplicantRepo,
	ur *repo.UserRepo,
) *CMSDashboardHandler {
	return &CMSDashboardHandler{nr, vr, ar, ur}
}

// GET /cms/dashboard
func (h *CMSDashboardHandler) Index(c *gin.Context) {
	newsCounts := h.newsRepo.CountByStatus()

	_, totalVolunteers, _ := h.volunteerRepo.FindAll(domain.ListOptions{
		Page: 1, Limit: 1, Status: "active",
	})
	totalApplicants := h.applicantRepo.CountTotal()

	mustChangePassword := false
	uid := sessionUserID(c)
	if user, err := h.userRepo.FindByID(uid); err == nil {
		mustChangePassword = user.MustChangePassword
	}

	c.HTML(http.StatusOK, "cms/dashboard.html", withUserCtx(c, gin.H{
		"title":                "Dashboard",
		"active_menu":          "dashboard",
		"total_published":      newsCounts["published"],
		"total_draft":          newsCounts["draft"],
		"total_pending":        newsCounts["pending_review"],
		"total_volunteers":     totalVolunteers,
		"total_applicants":     totalApplicants,
		"must_change_password": mustChangePassword,
	}))
}
