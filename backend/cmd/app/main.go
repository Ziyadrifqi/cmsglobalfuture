package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/yayasan/cms/config"
	"github.com/yayasan/cms/internal/domain"
	"github.com/yayasan/cms/internal/handler"
	"github.com/yayasan/cms/internal/mailer"
	"github.com/yayasan/cms/internal/middleware"
	"github.com/yayasan/cms/internal/repo"
	"github.com/yayasan/cms/migrations"
)

func main() {
	cfg := config.Load()
	db := config.NewDB(cfg)
	migrations.Run(db)

	mailerSvc := mailer.New(mailer.Config{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
		FromName: cfg.SMTPFromName,
	})

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	store := cookie.NewStore([]byte(cfg.SessionSecret))
	r.Use(sessions.Sessions("cms_session", store))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(cfg.AllowedOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Static("/static", "./static")

	// SetFuncMap HARUS sebelum LoadHTMLFiles
	r.SetFuncMap(template.FuncMap{
		"deref": func(p *uint) uint {
			if p == nil {
				return 0
			}
			return *p
		},
		"hasSuffix": strings.HasSuffix,
		// mul: kali bilangan integer, dipakai untuk animation-delay di template
		"mul": func(a, b int) int { return a * b },
	})

	r.LoadHTMLFiles(
		"templates/layouts/cms.html",
		"templates/layouts/sidebar.html",
		"templates/layouts/sidebar_footer.html",
		"templates/layouts/page_loader.html",
		"templates/layouts/sidebar_toggle.html",
		"templates/auth/login.html",
		"templates/cms/dashboard.html",
		"templates/cms/news_list.html",
		"templates/cms/news_form.html",
		"templates/cms/volunteer_list.html",
		"templates/cms/volunteer_form.html",
		"templates/cms/applicant_list.html",
		"templates/cms/gallery_list.html",
		"templates/cms/gallery_form.html",
		"templates/cms/gallery_category_list.html",
		"templates/cms/banner_list.html",
		"templates/cms/banner_form.html",
		"templates/cms/forbidden.html",
		"templates/cms/user_list.html",
		"templates/cms/user_form.html",
		"templates/cms/profile.html",
		"templates/cms/page_settings_index.html",
		"templates/cms/page_settings_form.html",
	)

	// ── Repo ──────────────────────────────────────────────────────────────────
	userRepo := repo.NewUserRepo(db)
	newsRepo := repo.NewNewsRepo(db)
	volunteerRepo := repo.NewVolunteerRepo(db)
	applicantRepo := repo.NewApplicantRepo(db)
	pageRepo := repo.NewPageRepo(db)
	bannerRepo := repo.NewBannerRepo(db)
	galleryRepo := repo.NewGalleryRepo(db)
	galleryCatRepo := repo.NewGalleryCategoryRepo(db)
	settingRepo := repo.NewSettingRepo(db)

	// ── Handler ───────────────────────────────────────────────────────────────
	authH := handler.NewCMSAuthHandler(userRepo)
	dashH := handler.NewCMSDashboardHandler(newsRepo, volunteerRepo, applicantRepo, userRepo)
	newsH := handler.NewCMSNewsHandler(newsRepo)
	volunteerH := handler.NewCMSVolunteerHandler(volunteerRepo, applicantRepo)
	galleryH := handler.NewCMSGalleryHandler(galleryRepo, galleryCatRepo)
	userH := handler.NewCMSUserHandler(userRepo, mailerSvc, cfg.AppBaseURL)
	bannerH := handler.NewCMSBannerHandler(bannerRepo)
	profileH := handler.NewCMSProfileHandler(userRepo)
	pageSettingsH := handler.NewCMSPageSettingsHandler(pageRepo, settingRepo)
	apiH := handler.NewAPIHandler(newsRepo, volunteerRepo, applicantRepo, pageRepo, bannerRepo, galleryRepo, galleryCatRepo, settingRepo)

	// ── Routes ────────────────────────────────────────────────────────────────

	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/cms/login") })

	// Auth
	r.GET("/cms/login", authH.LoginPage)
	r.POST("/cms/login", authH.Login)
	r.POST("/cms/logout", authH.Logout)

	// CMS (butuh login)
	cms := r.Group("/cms")
	cms.Use(middleware.RequireLogin())
	cms.Use(middleware.SetUserToContext())
	cms.Use(middleware.IdleTimeout())
	{
		cms.GET("/dashboard", dashH.Index)

		cms.GET("/profile", profileH.Index)
		cms.POST("/profile", profileH.Update)
		cms.POST("/profile/password", profileH.UpdatePassword)

		ng := cms.Group("/news")
		ng.Use(middleware.RequireRole(domain.RoleSuperAdmin, domain.RoleContentEditor, domain.RoleReviewer))
		{
			ng.GET("", newsH.Index)
			ng.GET("/create", newsH.CreatePage)
			ng.POST("/create", newsH.Create)
			ng.GET("/:id/edit", newsH.EditPage)
			ng.POST("/:id/edit", newsH.Update)
			ng.POST("/:id/submit-review", newsH.SubmitReview)
			ng.POST("/:id/approve", middleware.RequireRole(domain.RoleSuperAdmin, domain.RoleReviewer), newsH.Approve)
			ng.POST("/:id/reject", middleware.RequireRole(domain.RoleSuperAdmin, domain.RoleReviewer), newsH.Reject)
			ng.POST("/:id/delete", middleware.RequireRole(domain.RoleSuperAdmin), newsH.Delete)
		}

		bg := cms.Group("/banners")
		bg.Use(middleware.RequireRole(domain.RoleSuperAdmin, domain.RoleContentEditor))
		{
			bg.GET("", bannerH.Index)
			bg.GET("/create", bannerH.CreatePage)
			bg.POST("/create", bannerH.Create)
			bg.GET("/:id/edit", bannerH.EditPage)
			bg.POST("/:id/edit", bannerH.Update)
			bg.POST("/:id/toggle", bannerH.Toggle)
			bg.POST("/:id/delete", bannerH.Delete)
		}

		// ── Pengaturan Halaman: gabungan Page content + Site Settings per halaman
		psGroup := cms.Group("/page-settings")
		psGroup.Use(middleware.RequireRole(domain.RoleSuperAdmin, domain.RoleContentEditor))
		{
			psGroup.GET("", pageSettingsH.Index)
			psGroup.GET("/:page/edit", pageSettingsH.Edit)
			psGroup.POST("/:page/save", pageSettingsH.Save)
		}

		vg := cms.Group("/volunteers")
		vg.Use(middleware.RequireRole(domain.RoleSuperAdmin, domain.RoleHR))
		{
			vg.GET("", volunteerH.Index)
			vg.GET("/create", volunteerH.CreatePage)
			vg.POST("/create", volunteerH.Create)
			vg.GET("/:id/edit", volunteerH.EditPage)
			vg.POST("/:id/edit", volunteerH.Update)
			vg.POST("/:id/close", volunteerH.Close)
			vg.POST("/:id/delete", volunteerH.Delete)
			vg.GET("/:id/applicants", volunteerH.Applicants)
			vg.POST("/applicants/:id/status", volunteerH.UpdateApplicantStatus)
			vg.GET("/applicants/:id/cv", volunteerH.DownloadCV)
		}

		gg := cms.Group("/gallery")
		gg.Use(middleware.RequireRole(domain.RoleSuperAdmin, domain.RoleContentEditor))
		{
			gg.GET("", galleryH.Index)
			gg.GET("/create", galleryH.CreatePage)
			gg.POST("/create", galleryH.Create)
			gg.GET("/:id/edit", galleryH.EditPage)
			gg.POST("/:id/edit", galleryH.Update)
			gg.POST("/:id/toggle", galleryH.Toggle)
			gg.POST("/:id/delete", galleryH.Delete)
			gg.GET("/categories", galleryH.CategoryIndex)
			gg.POST("/categories/create", galleryH.CategoryCreate)
			gg.POST("/categories/:id/edit", galleryH.CategoryUpdate)
			gg.POST("/categories/:id/delete", galleryH.CategoryDelete)
		}

		ug := cms.Group("/users")
		ug.Use(middleware.RequireRole(domain.RoleSuperAdmin))
		{
			ug.GET("", userH.Index)
			ug.GET("/create", userH.CreatePage)
			ug.POST("/create", userH.Create)
			ug.GET("/:id/edit", userH.EditPage)
			ug.POST("/:id/edit", userH.Update)
			ug.POST("/:id/reset-password", userH.ResetPassword)
			ug.POST("/:id/delete", userH.Delete)
		}
	}

	// API untuk portal React
	api := r.Group("/api/v1")
	{
		api.GET("/home", apiH.Home)
		api.GET("/banners", apiH.Banners)
		api.GET("/settings", apiH.Settings)
		api.GET("/news", apiH.NewsList)
		api.GET("/news/:slug", apiH.NewsDetail)
		api.GET("/volunteers", apiH.VolunteerList)
		api.GET("/volunteers/:slug", apiH.VolunteerDetail)
		api.POST("/volunteers/:slug/apply", apiH.Apply)
		api.GET("/gallery", apiH.GalleryList)
		api.GET("/gallery/categories", apiH.GalleryCategories)
		api.GET("/pages/:slug", apiH.PageDetail)
	}

	addr := ":" + cfg.AppPort
	log.Printf("Server jalan di http://localhost%s", addr)
	log.Printf("CMS Panel : http://localhost%s/cms/login", addr)
	r.Run(addr)
}
