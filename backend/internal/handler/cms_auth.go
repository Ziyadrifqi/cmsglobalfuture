package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/yayasan/cms/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type CMSAuthHandler struct {
	userRepo *repo.UserRepo
}

func NewCMSAuthHandler(userRepo *repo.UserRepo) *CMSAuthHandler {
	return &CMSAuthHandler{userRepo}
}

// GET /cms/login
func (h *CMSAuthHandler) LoginPage(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("user_id") != nil {
		c.Redirect(http.StatusFound, "/cms/dashboard")
		return
	}
	// Tampilkan pesan khusus jika di-redirect karena idle timeout
	c.HTML(http.StatusOK, "auth/login.html", gin.H{
		"title":  "Login CMS",
		"reason": c.Query("reason"), // "idle" → pesan timeout
	})
}

// POST /cms/login
func (h *CMSAuthHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := h.userRepo.FindByEmail(email)
	if err != nil || !user.IsActive {
		c.HTML(http.StatusOK, "auth/login.html", gin.H{
			"title": "Login CMS",
			"error": "Email atau password salah",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.HTML(http.StatusOK, "auth/login.html", gin.H{
			"title": "Login CMS",
			"error": "Email atau password salah",
		})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_name", user.Name)
	session.Set("user_role", user.Role.Name)
	session.Set("user_role_display", user.Role.DisplayName)
	session.Set("user_avatar", user.Avatar)
	session.Save()

	c.Redirect(http.StatusFound, "/cms/dashboard")
}

// POST /cms/logout
func (h *CMSAuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/cms/login")
}
