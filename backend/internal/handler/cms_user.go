package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yayasan/cms/internal/domain"
	"github.com/yayasan/cms/internal/mailer"
	"github.com/yayasan/cms/internal/repo"
	"github.com/yayasan/cms/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type CMSUserHandler struct {
	userRepo   *repo.UserRepo
	mailer     *mailer.Mailer
	appBaseURL string
}

func NewCMSUserHandler(ur *repo.UserRepo, m *mailer.Mailer, appBaseURL string) *CMSUserHandler {
	return &CMSUserHandler{ur, m, appBaseURL}
}

// GET /cms/users
func (h *CMSUserHandler) Index(c *gin.Context) {
	users, _ := h.userRepo.FindAll()
	c.HTML(http.StatusOK, "cms/user_list.html", withUserCtx(c, gin.H{
		"title":       "Manajemen Pengguna",
		"active_menu": "users",
		"users":       users,
		"flash":       c.Query("flash"),
	}))
}

// GET /cms/users/create
func (h *CMSUserHandler) CreatePage(c *gin.Context) {
	roles, _ := h.userRepo.FindAllRoles()
	c.HTML(http.StatusOK, "cms/user_form.html", withUserCtx(c, gin.H{
		"title":       "Tambah Pengguna",
		"active_menu": "users",
		"roles":       roles,
		"is_edit":     false,
	}))
}

// POST /cms/users/create
//
// BARU: pengiriman email sekarang SYNCHRONOUS (bukan lagi `go func(...)`).
// Alasannya: dengan goroutine async, handler langsung redirect dengan
// flash=created sebelum tahu hasil kirim email sebenarnya — jadi kalau
// SMTP salah konfigurasi atau gagal connect, admin tidak pernah tahu dari
// UI, cuma kelihatan di log server. Mengirim email lewat SMTP biasanya
// makan waktu < 1-2 detik, jadi delay ini biasanya tidak terasa berat
// untuk aksi "buat pengguna" yang memang tidak sering dilakukan.
//
// Akun TETAP dibuat meskipun email gagal terkirim (supaya admin tidak
// kehilangan kerjaan kalau SMTP lagi bermasalah) — tapi flash message
// akan membedakan dengan jelas: "dibuat & email terkirim" vs
// "dibuat TAPI email gagal terkirim, hubungi user manual".
func (h *CMSUserHandler) Create(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.PostForm("role_id"))
	name := c.PostForm("name")
	email := c.PostForm("email")

	roles, _ := h.userRepo.FindAllRoles()
	var selectedRole *domain.Role
	for i := range roles {
		if roles[i].ID == uint(roleID) {
			selectedRole = &roles[i]
			break
		}
	}
	if selectedRole == nil {
		c.HTML(http.StatusOK, "cms/user_form.html", withUserCtx(c, gin.H{
			"title":       "Tambah Pengguna",
			"active_menu": "users",
			"roles":       roles,
			"is_edit":     false,
			"error_info":  "Role yang dipilih tidak valid",
		}))
		return
	}

	// Generate password acak unik — TIDAK sama antar user, TIDAK ada di
	// konstanta kode manapun.
	generatedPassword, err := util.GenerateSecurePassword(12)
	if err != nil {
		c.HTML(http.StatusOK, "cms/user_form.html", withUserCtx(c, gin.H{
			"title":       "Tambah Pengguna",
			"active_menu": "users",
			"roles":       roles,
			"is_edit":     false,
			"error_info":  "Gagal membuat password acak, coba lagi",
		}))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(generatedPassword), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusOK, "cms/user_form.html", withUserCtx(c, gin.H{
			"title":       "Tambah Pengguna",
			"active_menu": "users",
			"roles":       roles,
			"is_edit":     false,
			"error_info":  "Gagal membuat akun, coba lagi",
		}))
		return
	}

	newUser := &domain.User{
		Name:               name,
		Email:              email,
		Password:           string(hash), // hanya hash yang disimpan, plaintext tidak pernah masuk DB
		IsActive:           true,
		RoleID:             uint(roleID),
		MustChangePassword: true,
	}

	if err := h.userRepo.Create(newUser); err != nil {
		c.HTML(http.StatusOK, "cms/user_form.html", withUserCtx(c, gin.H{
			"title":       "Tambah Pengguna",
			"active_menu": "users",
			"roles":       roles,
			"is_edit":     false,
			"error_info":  "Gagal menyimpan pengguna. Pastikan email belum dipakai.",
		}))
		return
	}

	// Kirim email berisi password plaintext SEKALI SAJA, secara synchronous
	// (ditunggu hasilnya) — supaya kita bisa redirect dengan flash yang
	// sesuai dengan apa yang BENAR-BENAR terjadi.
	loginURL := h.appBaseURL + "/cms/login"
	sendErr := h.mailer.SendNewAccountEmail(email, name, generatedPassword, selectedRole.DisplayName, loginURL)

	c.Redirect(http.StatusFound, "/cms/users?flash="+emailFlashFor(sendErr, "created"))
}

// GET /cms/users/:id/edit
func (h *CMSUserHandler) EditPage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/users")
		return
	}
	roles, _ := h.userRepo.FindAllRoles()
	c.HTML(http.StatusOK, "cms/user_form.html", withUserCtx(c, gin.H{
		"title":       "Edit Pengguna",
		"active_menu": "users",
		"user":        user,
		"roles":       roles,
		"is_edit":     true,
	}))
}

// POST /cms/users/:id/edit
func (h *CMSUserHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/users")
		return
	}
	roleID, _ := strconv.Atoi(c.PostForm("role_id"))
	user.Name = c.PostForm("name")
	user.Email = c.PostForm("email")
	user.RoleID = uint(roleID)
	user.IsActive = c.PostForm("is_active") == "1"

	if pw := c.PostForm("password"); pw != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		user.Password = string(hash)
		user.MustChangePassword = true
	}
	h.userRepo.Update(user)
	c.Redirect(http.StatusFound, "/cms/users?flash=updated")
}

// POST /cms/users/:id/reset-password
//
// BARU: sama seperti Create() di atas — pengiriman email sekarang
// SYNCHRONOUS, dan flash message membedakan dengan jelas antara:
//   - password_reset            → reset sukses, email terkirim
//   - reset_email_not_configured → reset sukses, TAPI SMTP belum diisi di .env
//   - reset_email_failed         → reset sukses, TAPI email gagal kirim (cek log untuk detail SMTP error)
//   - reset_failed               → reset GAGAL total (password tidak diganti)
//
// Password tetap diganti di DB meskipun email gagal — supaya tindakan
// reset tidak hilang sia-sia. Tapi sekarang admin TAHU kalau perlu kasih
// tahu user secara manual karena email tidak nyampe.
func (h *CMSUserHandler) ResetPassword(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/users")
		return
	}

	newPassword, err := util.GenerateSecurePassword(12)
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/users?flash=reset_failed")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Redirect(http.StatusFound, "/cms/users?flash=reset_failed")
		return
	}

	user.Password = string(hash)
	user.MustChangePassword = true
	if err := h.userRepo.Update(user); err != nil {
		c.Redirect(http.StatusFound, "/cms/users?flash=reset_failed")
		return
	}

	loginURL := h.appBaseURL + "/cms/login"
	sendErr := h.mailer.SendPasswordResetEmail(user.Email, user.Name, newPassword, loginURL)

	c.Redirect(http.StatusFound, "/cms/users?flash="+emailFlashFor(sendErr, "password_reset"))
}

// POST /cms/users/:id/delete
func (h *CMSUserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.userRepo.Delete(uint(id))
	c.Redirect(http.StatusFound, "/cms/users?flash=deleted")
}

// ── helper ────────────────────────────────────────────────────────────────────

// emailFlashFor menentukan kode flash berdasarkan hasil pengiriman email.
// successFlash adalah kode flash yang dipakai kalau email berhasil terkirim
// (misal "created" atau "password_reset") — supaya pesan sukses lama tetap
// dipakai apa adanya, dan cuma menambah varian baru untuk kasus gagal.
func emailFlashFor(sendErr error, successFlash string) string {
	if sendErr == nil {
		return successFlash
	}
	if errors.Is(sendErr, mailer.ErrSMTPNotConfigured) {
		return successFlash + "_email_not_configured"
	}
	return successFlash + "_email_failed"
}
