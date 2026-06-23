package handler

import "github.com/gin-gonic/gin"

// withUserCtx menggabungkan data user dari Gin context ke dalam gin.H
// yang dikirim ke template. Versi ini AMAN: pakai c.Get() bukan c.MustGet()
// sehingga tidak panic kalau key tidak ada (misal handler dipakai di luar
// middleware SetUserToContext).
func withUserCtx(c *gin.Context, data gin.H) gin.H {
	if data == nil {
		data = gin.H{}
	}
	keys := []string{
		"session_user_id",
		"session_user_name",
		"session_user_role",
		"session_user_role_display",
		"session_user_avatar",
	}
	for _, k := range keys {
		if val, exists := c.Get(k); exists {
			data[k] = val
		}
	}
	return data
}

// sessionUserID mengambil user ID dari Gin context dengan aman.
// Mengembalikan 0 jika key tidak ada atau tipe tidak cocok — tidak panic.
func sessionUserID(c *gin.Context) uint {
	val, exists := c.Get("session_user_id")
	if !exists {
		return 0
	}
	switch v := val.(type) {
	case uint:
		return v
	case int:
		return uint(v)
	case uint32:
		return uint(v)
	case uint64:
		return uint(v)
	case float64: // gob decode kadang hasilkan float64
		return uint(v)
	}
	return 0
}
