package handler

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// withUserCtx menggabungkan data user dari Gin context ke dalam gin.H.
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
	case float64:
		return uint(v)
	}
	return 0
}

// removeFile menghapus file upload dari disk secara aman.
//
// urlPath adalah path URL seperti yang disimpan di database,
// contoh: "/static/uploads/banners/foto.jpg"
//
// Hanya file di dalam folder static/uploads/ yang boleh dihapus —
// ada safety check supaya tidak bisa hapus file di luar folder upload.
// Error "file tidak ditemukan" diabaikan secara diam-diam (tidak perlu
// dikhawatirkan kalau file memang sudah tidak ada).
func removeFile(urlPath string) {
	if urlPath == "" {
		return
	}
	// Konversi URL path → filesystem path: "/static/..." → "static/..."
	fspath := strings.TrimPrefix(urlPath, "/")

	// Guard: hanya hapus file di dalam static/uploads/
	if !strings.HasPrefix(fspath, "static/uploads/") {
		log.Printf("removeFile: ditolak — path di luar uploads: %s", fspath)
		return
	}

	if err := os.Remove(fspath); err != nil && !os.IsNotExist(err) {
		log.Printf("removeFile: gagal hapus %s: %v", fspath, err)
	}
}
