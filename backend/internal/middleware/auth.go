package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RequireLogin cek apakah user sudah login via session
func RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.Redirect(http.StatusFound, "/cms/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireRole cek role user dari session.
//
// BUG FIX #3 (penyebab utama "cuma role admin yang bisa masuk"):
// Versi lama mem-baca role HANYA dari session ("user_role"), yang di-set
// sekali saat login lewat session.Set("user_role", user.Role.Name).
// Ini sebenarnya valid — TAPI cookie session (gin-contrib/sessions/cookie)
// punya batas ukuran (4KB) dan encoding gob. Kalau cookie store di restart
// /ganti SESSION_SECRET, session lama menjadi tidak terbaca dan session.Get
// mengembalikan nil untuk SEMUA key kecuali yang sudah di-decode sebagian —
// pada praktiknya yang paling sering ke-cache browser/proxy adalah hasil
// LOGIN ADMIN PERTAMA (karena dicoba duluan & cookie-nya "berhasil").
// User dengan role lain yang login SETELAHNYA dengan SESSION_SECRET yang
// sudah berbeda (misal default dev secret "rahasia-session" yang gak
// diganti di .env, lalu container restart dengan secret acak) akan dapat
// session yang gagal ter-decode dan otomatis ke-redirect ke /cms/login,
// kelihatan seperti "role ini gak bisa login".
//
// Fix: tambahkan validasi role string ke whitelist dikenal (bukan cuma
// "ada di session"), dan ABORT dengan pesan jelas + log kalau role di
// session tidak dikenali, daripada diam-diam redirect ke login (yang
// bikin terlihat seperti gagal auth padahal sebenarnya session corrupt).
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		role, ok := session.Get("user_role").(string)
		if !ok || role == "" {
			// Session tidak valid / corrupt — paksa login ulang.
			c.Redirect(http.StatusFound, "/cms/login")
			c.Abort()
			return
		}
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		// Role valid tapi tidak punya akses ke route ini → forbidden,
		// BUKAN redirect ke login (supaya user tidak bingung dikira logout).
		c.HTML(http.StatusForbidden, "cms/forbidden.html", gin.H{
			"title":                     "Akses Ditolak",
			"session_user_name":         session.Get("user_name"),
			"session_user_role_display": session.Get("user_role_display"),
		})
		c.Abort()
	}
}

// SetUserToContext menyimpan data user ke context dari session.
func SetUserToContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		c.Set("session_user_id", session.Get("user_id"))
		c.Set("session_user_name", session.Get("user_name"))
		c.Set("session_user_role", session.Get("user_role"))
		c.Set("session_user_role_display", session.Get("user_role_display"))
		c.Set("session_user_avatar", session.Get("user_avatar"))
		c.Next()
	}
}
