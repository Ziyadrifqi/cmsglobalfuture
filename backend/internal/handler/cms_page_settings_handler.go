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

// pageDef mendefinisikan satu "halaman" di menu Pengaturan Halaman.
// Sebuah halaman bisa punya konten panjang dari tabel pages (PageSlug tidak
// kosong) dan/atau kumpulan Setting key yang spesifik untuk halaman itu
// (SettingKeys). Form edit akan menampilkan keduanya sekaligus dalam satu form.
type pageDef struct {
	Name        string
	Slug        string // slug untuk URL /cms/page-settings/:slug/edit
	Icon        string
	Description string
	PageSlug    string   // slug di tabel domain.Page — kosong kalau tidak ada
	SettingKeys []string // key dari tabel site_settings yang relevan
}

// pageDefinitions adalah urutan & definisi halaman yang tampil di menu.
// Perubahan urutan atau penambahan halaman baru cukup di sini.
var pageDefinitions = []pageDef{
	{
		Name:        "Beranda",
		Slug:        "beranda",
		Icon:        "🏠",
		Description: "Hero, statistik, program, dampak, dan CTA halaman utama",
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
		Description: "Teks deskripsi organisasi, tim pengurus, nilai-nilai, dan linimasa",
		PageSlug:    "about",
		SettingKeys: []string{
			"about_team_json",
			"about_values_json",
			"about_milestones_json",
		},
	},
	{
		Name:        "Visi & Misi",
		Slug:        "visi-misi",
		Icon:        "🔭",
		Description: "Teks visi dan misi organisasi",
		PageSlug:    "vision-mission",
		SettingKeys: []string{},
	},
	{
		Name:        "Kontak",
		Slug:        "kontak",
		Icon:        "📞",
		Description: "Informasi kontak, jam operasional, alamat, dan akun sosial media",
		PageSlug:    "contact",
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
		Description: "Konten di halaman daftar rekrutmen relawan",
		PageSlug:    "",
		SettingKeys: []string{
			"volunteer_why_join_json",
		},
	},
	{
		Name:        "Footer & Identitas Situs",
		Slug:        "footer-identitas",
		Icon:        "🏷️",
		Description: "Nama situs, logo, deskripsi footer, daftar program, dan akun sosial media footer",
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

// findPageDef mencari definisi halaman berdasarkan slug URL.
func findPageDef(slug string) *pageDef {
	for i := range pageDefinitions {
		if pageDefinitions[i].Slug == slug {
			return &pageDefinitions[i]
		}
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────

// CMSPageSettingsHandler menangani menu Pengaturan Halaman yang menggabungkan
// tabel pages (teks panjang) dan tabel site_settings (konten kecil dinamis).
type CMSPageSettingsHandler struct {
	pageRepo    *repo.PageRepo
	settingRepo *repo.SettingRepo
}

func NewCMSPageSettingsHandler(pr *repo.PageRepo, sr *repo.SettingRepo) *CMSPageSettingsHandler {
	return &CMSPageSettingsHandler{pr, sr}
}

// GET /cms/page-settings — daftar kartu halaman
func (h *CMSPageSettingsHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "cms/page_settings_index.html", withUserCtx(c, gin.H{
		"title":       "Pengaturan Halaman",
		"active_menu": "page_settings",
		"pages":       pageDefinitions,
		"flash":       c.Query("flash"),
	}))
}

// GET /cms/page-settings/:page/edit — form edit satu halaman
func (h *CMSPageSettingsHandler) Edit(c *gin.Context) {
	def := findPageDef(c.Param("page"))
	if def == nil {
		c.Redirect(http.StatusFound, "/cms/page-settings")
		return
	}

	// Ambil konten teks panjang dari tabel pages (kalau ada)
	var pageContent *domain.Page
	if def.PageSlug != "" {
		if p, err := h.pageRepo.FindBySlug(def.PageSlug); err == nil {
			pageContent = p
		}
	}

	// Ambil setting yang relevan untuk halaman ini
	allSettings, _ := h.settingRepo.FindAll()
	keySet := map[string]bool{}
	for _, k := range def.SettingKeys {
		keySet[k] = true
	}
	var fields []domain.SiteSetting
	for _, s := range allSettings {
		if keySet[s.Key] {
			fields = append(fields, s)
		}
	}

	// Urutkan fields sesuai urutan SettingKeys (supaya form rapi)
	ordered := make([]domain.SiteSetting, 0, len(fields))
	for _, k := range def.SettingKeys {
		for _, f := range fields {
			if f.Key == k {
				ordered = append(ordered, f)
				break
			}
		}
	}

	c.HTML(http.StatusOK, "cms/page_settings_form.html", withUserCtx(c, gin.H{
		"title":        def.Name,
		"active_menu":  "page_settings",
		"page_def":     def,
		"page_content": pageContent, // bisa nil kalau tidak ada PageSlug
		"fields":       ordered,
		"flash":        c.Query("flash"),
	}))
}

// POST /cms/page-settings/:page/save — simpan semua perubahan halaman ini
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

	// Simpan settings yang relevan untuk halaman ini
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
				// kalau tidak upload file baru, nilai lama tetap dipertahankan
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
