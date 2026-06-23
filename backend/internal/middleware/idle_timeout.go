package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	sessionLastActiveKey = "last_active"
	idleTimeout          = 60 * time.Minute
)

// IdleTimeout logout otomatis jika tidak ada aktivitas selama 1 jam.
// Middleware ini harus dipasang SETELAH RequireLogin() pada group /cms.
func IdleTimeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		raw := session.Get(sessionLastActiveKey)
		if raw != nil {
			var lastActive time.Time
			switch v := raw.(type) {
			case int64:
				lastActive = time.Unix(v, 0)
			case time.Time:
				lastActive = v
			}
			if !lastActive.IsZero() && time.Since(lastActive) > idleTimeout {
				session.Clear()
				session.Save()
				c.Redirect(http.StatusFound, "/cms/login?reason=idle")
				c.Abort()
				return
			}
		}

		// Perbarui timestamp aktivitas terakhir
		session.Set(sessionLastActiveKey, time.Now().Unix())
		session.Save()

		c.Next()
	}
}
