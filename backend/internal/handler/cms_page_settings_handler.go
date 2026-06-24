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

type pageDef struct {
	Name        string
	Slug        string
	Icon        string
	Description string
	PageSlug    string
	SettingKeys []string
}

// pageDefinitions mendefinisikan urutan dan isi tiap kartu di menu
// "Pengaturan Halaman". Sesuaikan dengan halaman yang ada di frontend React.
var pageDefinitions = []pageDef{
	{
		Name:        "Beranda",
		Slug:        "beranda",
		Icon:        "🏠",
		Description: "Hero, statistik, program unggulan, dampak nyata, dan CTA",
		PageSlug:    "",
		SettingKeys: []string{
			"hero_badge_text",
			"hero_title_main",
			"hero_title_highlight",
			"hero_subtitle",
			"hero_cta_primary_text",
			"hero_cta_secondary_text",
			"home_stats_json",
			"home_programs_json",
			"home_impacts_json",
			"home_cta_title",
			"home_cta_subtitle",
		},
	},
	{
		Name:        "Tentang Kami",
		Slug:        "tentang-kami",
		Icon:        "👥",
		Description: "Deskripsi organisasi, visi & misi, linimasa, nilai-nilai, dan tim pengurus",
		PageSlug:    "about",
		SettingKeys: []string{
			"about_vision_text",
			"about_mission_items",
			"about_milestones_json",
			"about_values_json",
			"about_team_json",
		},
	},
	{
		Name:        "Kontak",
		Slug:        "kontak",
		Icon:        "📞",
		Description: "Email, telepon, alamat, jam operasional, dan akun sosial media",
		PageSlug:    "",
		SettingKeys: []string{
			"contact_email_1",
			"contact_email_2",
			"contact_phone",
			"contact_phone_note",
			"contact_address_1",
			"contact_address_2",
			"contact_hours_1",
			"contact_hours_2",
			"contact_map_label",
			"contact_social_json",
		},
	},
	{
		Name:        "Halaman Relawan",
		Slug:        "relawan",
		Icon:        "🤝",
		Description: "Alasan bergabung menjadi relawan yang tampil di halaman daftar rekrutmen",
		PageSlug:    "",
		SettingKeys: []string{
			"volunteer_why_join_json",
		},
	},
	{
		Name:        "Footer & Identitas Situs",
		Slug:        "footer-identitas",
		Icon:        "🏷️",
		Description: "Nama situs, logo, deskripsi footer, daftar program, dan ikon sosial media",
		PageSlug:    "",
		SettingKeys: []string{
			"site_name",
			"site_name_sub",
			"site_logo_image",
			"footer_description",
			"footer_programs_json",
			"footer_social_json",
		},
	},
}

func findPageDef(slug string) *pageDef {
	for i := range pageDefinitions {
		if pageDefinitions[i].Slug == slug {
			return &pageDefinitions[i]
		}
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────

type CMSPageSettingsHandler struct {
	pageRepo    *repo.PageRepo
	settingRepo *repo.SettingRepo
}

func NewCMSPageSettingsHandler(pr *repo.PageRepo, sr *repo.SettingRepo) *CMSPageSettingsHandler {
	return &CMSPageSettingsHandler{pr, sr}
}

// GET /cms/page-settings
func (h *CMSPageSettingsHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "cms/page_settings_index.html", withUserCtx(c, gin.H{
		"title":       "Pengaturan Halaman",
		"active_menu": "page_settings",
		"pages":       pageDefinitions,
		"flash":       c.Query("flash"),
	}))
}

// GET /cms/page-settings/:page/edit
func (h *CMSPageSettingsHandler) Edit(c *gin.Context) {
	def := findPageDef(c.Param("page"))
	if def == nil {
		c.Redirect(http.StatusFound, "/cms/page-settings")
		return
	}

	var pageContent *domain.Page
	if def.PageSlug != "" {
		if p, err := h.pageRepo.FindBySlug(def.PageSlug); err == nil {
			pageContent = p
		}
	}

	allSettings, _ := h.settingRepo.FindAll()
	keySet := map[string]bool{}
	for _, k := range def.SettingKeys {
		keySet[k] = true
	}

	// Urutkan fields sesuai urutan SettingKeys supaya form rapi
	fieldMap := map[string]domain.SiteSetting{}
	for _, s := range allSettings {
		if keySet[s.Key] {
			fieldMap[s.Key] = s
		}
	}
	ordered := make([]domain.SiteSetting, 0, len(def.SettingKeys))
	for _, k := range def.SettingKeys {
		if f, ok := fieldMap[k]; ok {
			ordered = append(ordered, f)
		}
	}

	c.HTML(http.StatusOK, "cms/page_settings_form.html", withUserCtx(c, gin.H{
		"title":        def.Name,
		"active_menu":  "page_settings",
		"page_def":     def,
		"page_content": pageContent,
		"fields":       ordered,
		"flash":        c.Query("flash"),
	}))
}

// POST /cms/page-settings/:page/save
func (h *CMSPageSettingsHandler) Save(c *gin.Context) {
	def := findPageDef(c.Param("page"))
	if def == nil {
		c.Redirect(http.StatusFound, "/cms/page-settings")
		return
	}

	// Simpan konten teks panjang (Page) kalau ada
	if def.PageSlug != "" {
		if p, err := h.pageRepo.FindBySlug(def.PageSlug); err == nil {
			if title := c.PostForm("page_title"); title != "" {
				p.Title = title
			}
			p.Content = c.PostForm("page_content")
			p.MetaTitle = c.PostForm("page_meta_title")
			p.MetaDescription = c.PostForm("page_meta_description")
			h.pageRepo.Update(p)
		}
	}

	// Simpan settings yang relevan
	if len(def.SettingKeys) > 0 {
		allSettings, _ := h.settingRepo.FindAll()
		keySet := map[string]bool{}
		for _, k := range def.SettingKeys {
			keySet[k] = true
		}
		for i := range allSettings {
			s := &allSettings[i]
			if !keySet[s.Key] {
				continue
			}
			if s.Type == domain.SettingTypeImage {
				file, err := c.FormFile("image[" + s.Key + "]")
				if err == nil {
					ext := filepath.Ext(file.Filename)
					filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + s.Key + ext
					dst := "static/uploads/settings/" + filename
					if c.SaveUploadedFile(file, dst) == nil {
						s.Value = "/static/uploads/settings/" + filename
						h.settingRepo.Update(s)
					}
				}
			} else {
				if v, ok := c.GetPostForm("setting[" + s.Key + "]"); ok {
					s.Value = v
					h.settingRepo.Update(s)
				}
			}
		}
	}

	c.Redirect(http.StatusFound, "/cms/page-settings/"+def.Slug+"/edit?flash=updated")
}
