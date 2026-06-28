package handler

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/yayasan/cms/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type CMSProfileHandler struct {
	userRepo *repo.UserRepo
}

func NewCMSProfileHandler(ur *repo.UserRepo) *CMSProfileHandler {
	return &CMSProfileHandler{ur}
}

func (h *CMSProfileHandler) Index(c *gin.Context) {
	uid := sessionUserID(c)
	user, err := h.userRepo.FindByID(uid)
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/dashboard")
		return
	}
	c.HTML(http.StatusOK, "cms/profile.html", withUserCtx(c, gin.H{
		"title": "Profil Saya", "active_menu": "profile",
		"user": user, "flash": c.Query("flash"),
	}))
}

func (h *CMSProfileHandler) Update(c *gin.Context) {
	uid := sessionUserID(c)
	user, err := h.userRepo.FindByID(uid)
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/dashboard")
		return
	}

	name := c.PostForm("name")
	email := c.PostForm("email")

	if name == "" || email == "" {
		c.HTML(http.StatusOK, "cms/profile.html", withUserCtx(c, gin.H{
			"title": "Profil Saya", "active_menu": "profile",
			"user": user, "error_info": "Nama dan email wajib diisi",
		}))
		return
	}

	if existing, err := h.userRepo.FindByEmail(email); err == nil && existing.ID != user.ID {
		c.HTML(http.StatusOK, "cms/profile.html", withUserCtx(c, gin.H{
			"title": "Profil Saya", "active_menu": "profile",
			"user": user, "error_info": "Email sudah digunakan oleh pengguna lain",
		}))
		return
	}

	user.Name = name
	user.Email = email

	file, ferr := c.FormFile("avatar")
	if ferr == nil {
		oldAvatar := user.Avatar // simpan path lama
		ext := filepath.Ext(file.Filename)
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_avatar" + ext
		dst := "static/uploads/avatars/" + filename
		if c.SaveUploadedFile(file, dst) == nil {
			user.Avatar = "/static/uploads/avatars/" + filename
			removeFile(oldAvatar) // hapus foto profil lama dari disk
		}
	}

	if err := h.userRepo.Update(user); err != nil {
		c.HTML(http.StatusOK, "cms/profile.html", withUserCtx(c, gin.H{
			"title": "Profil Saya", "active_menu": "profile",
			"user": user, "error_info": "Gagal menyimpan perubahan, coba lagi",
		}))
		return
	}

	session := sessions.Default(c)
	session.Set("user_name", user.Name)
	session.Set("user_avatar", user.Avatar)
	session.Save()

	c.Redirect(http.StatusFound, "/cms/profile?flash=updated")
}

func (h *CMSProfileHandler) UpdatePassword(c *gin.Context) {
	uid := sessionUserID(c)
	user, err := h.userRepo.FindByID(uid)
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/dashboard")
		return
	}

	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	confirmPassword := c.PostForm("confirm_password")

	renderError := func(msg string) {
		c.HTML(http.StatusOK, "cms/profile.html", withUserCtx(c, gin.H{
			"title": "Profil Saya", "active_menu": "profile",
			"user": user, "error_password": msg,
		}))
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)) != nil {
		renderError("Password lama yang Anda masukkan salah")
		return
	}
	if len(newPassword) < 6 {
		renderError("Password baru minimal 6 karakter")
		return
	}
	if newPassword != confirmPassword {
		renderError("Konfirmasi password baru tidak cocok")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		renderError("Gagal mengenkripsi password baru, coba lagi")
		return
	}

	user.Password = string(hash)
	user.MustChangePassword = false
	if err := h.userRepo.Update(user); err != nil {
		renderError("Gagal menyimpan password baru, coba lagi")
		return
	}

	c.Redirect(http.StatusFound, "/cms/profile?flash=password_updated")
}
