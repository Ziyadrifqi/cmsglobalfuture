package migrations

import (
	"log"

	"github.com/yayasan/cms/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	log.Println("Menjalankan migrasi...")
	err := db.AutoMigrate(
		&domain.Role{},
		&domain.User{},
		&domain.NewsCategory{},
		&domain.News{},
		&domain.Volunteer{},
		&domain.Applicant{},
		&domain.Banner{},
		&domain.Page{},
		&domain.GalleryCategory{},
		&domain.GalleryItem{},
		&domain.SiteSetting{},
	)
	if err != nil {
		log.Fatalf("Migrasi gagal: %v", err)
	}
	log.Println("✓ Migrasi selesai")
	seedAll(db)
}

func seedAll(db *gorm.DB) {
	seedRoles(db)
	seedAdmin(db)
	seedPages(db)
	seedGalleryCategories(db)
	seedSettings(db)
}

func seedRoles(db *gorm.DB) {
	roleDefs := []domain.Role{
		{Name: domain.RoleSuperAdmin, DisplayName: "Super Admin"},
		{Name: domain.RoleContentEditor, DisplayName: "Content Editor"},
		{Name: domain.RoleHR, DisplayName: "HR Recruitment"},
		{Name: domain.RoleReviewer, DisplayName: "Reviewer"},
	}
	for _, def := range roleDefs {
		role := def
		if err := db.Where(domain.Role{Name: role.Name}).FirstOrCreate(&role).Error; err != nil {
			log.Fatalf("Gagal seed role %s: %v", role.Name, err)
		}
	}
	log.Println("✓ Role di-seed")
}

func seedAdmin(db *gorm.DB) {
	var role domain.Role
	if err := db.Where("name = ?", domain.RoleSuperAdmin).First(&role).Error; err != nil {
		log.Fatalf("Role super_admin tidak ditemukan: %v", err)
	}
	var existing domain.User
	if err := db.Where("email = ?", "admin@yayasan.local").First(&existing).Error; err == nil {
		if existing.RoleID != role.ID {
			existing.RoleID = role.ID
			db.Save(&existing)
		}
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
	db.Create(&domain.User{
		Name:     "Super Admin",
		Email:    "admin@yayasan.local",
		Password: string(hash),
		IsActive: true,
		RoleID:   role.ID,
	})
	log.Println("✓ Admin default: admin@yayasan.local / Admin@123")
}

func seedPages(db *gorm.DB) {
	pages := []domain.Page{
		{Slug: "about", Title: "Tentang Kami", Content: "Green Future Indonesia adalah organisasi lingkungan nirlaba yang didirikan pada tahun 2015 oleh sekumpulan aktivis muda yang prihatin terhadap kondisi lingkungan Indonesia.\n\nKami percaya bahwa perubahan lingkungan yang nyata hanya bisa terwujud melalui kolaborasi antara masyarakat, pemerintah, dan sektor swasta. Karena itu, kami membangun ekosistem relawan yang kuat — dari pelajar, mahasiswa, profesional, hingga komunitas lokal.\n\nHingga 2025, Green Future Indonesia telah menjangkau lebih dari 28 kota dengan 12.000+ relawan aktif."},
	}
	for _, def := range pages {
		p := def
		db.Where(domain.Page{Slug: p.Slug}).FirstOrCreate(&p)
	}
	log.Println("✓ Halaman statis di-seed")
}

func seedGalleryCategories(db *gorm.DB) {
	cats := []domain.GalleryCategory{
		{Name: "Penanaman Mangrove", Slug: "penanaman-mangrove", OrderNum: 1},
		{Name: "Bersih Pantai", Slug: "bersih-pantai", OrderNum: 2},
		{Name: "Edukasi", Slug: "edukasi", OrderNum: 3},
		{Name: "Workshop", Slug: "workshop", OrderNum: 4},
		{Name: "Event Nasional", Slug: "event-nasional", OrderNum: 5},
	}
	for _, def := range cats {
		cat := def
		db.Where(domain.GalleryCategory{Slug: cat.Slug}).FirstOrCreate(&cat)
	}
	log.Println("✓ Kategori galeri di-seed")
}

func seedSettings(db *gorm.DB) {
	defs := []domain.SiteSetting{
		// ── Identitas Situs ─────────────────────────────────────────────────
		{Key: "site_name", Label: "Nama Situs (baris atas)", Section: "Identitas Situs", Type: domain.SettingTypeText, OrderNum: 1, Value: "Green Future"},
		{Key: "site_name_sub", Label: "Nama Situs (baris bawah)", Section: "Identitas Situs", Type: domain.SettingTypeText, OrderNum: 2, Value: "Indonesia"},
		{Key: "site_logo_image", Label: "Logo (opsional — kosongkan untuk pakai ikon 🌿)", Section: "Identitas Situs", Type: domain.SettingTypeImage, OrderNum: 3, Value: ""},

		// ── Hero Beranda ─────────────────────────────────────────────────────
		{Key: "hero_badge_text", Label: "Teks Badge di atas judul", Section: "Hero Beranda", Type: domain.SettingTypeText, OrderNum: 1, Value: "12.000+ Relawan Aktif di 28 Kota"},
		{Key: "hero_title_main", Label: "Judul Hero (baris utama)", Section: "Hero Beranda", Type: domain.SettingTypeText, OrderNum: 2, Value: "Bersama Jaga Bumi"},
		{Key: "hero_title_highlight", Label: "Judul Hero (teks sorotan kuning)", Section: "Hero Beranda", Type: domain.SettingTypeText, OrderNum: 3, Value: "untuk Generasi Depan"},
		{Key: "hero_subtitle", Label: "Deskripsi / Subjudul Hero", Section: "Hero Beranda", Type: domain.SettingTypeText, OrderNum: 4, Value: "Green Future Indonesia adalah gerakan lingkungan yang mengajak semua orang untuk beraksi nyata — menanam, membersihkan, dan mendidik demi Indonesia yang lebih hijau."},
		{Key: "hero_cta_primary_text", Label: "Tombol Utama (teks)", Section: "Hero Beranda", Type: domain.SettingTypeText, OrderNum: 5, Value: "Jadi Relawan"},
		{Key: "hero_cta_secondary_text", Label: "Tombol Kedua (teks)", Section: "Hero Beranda", Type: domain.SettingTypeText, OrderNum: 6, Value: "Tentang Kami"},

		// ── Statistik Beranda ────────────────────────────────────────────────
		{Key: "home_stats_json", Label: "Kartu Statistik (ikon, angka, keterangan)", Section: "Statistik Beranda", Type: domain.SettingTypeJSON, OrderNum: 1,
			Value: `[{"icon":"🌿","value":"12K+","label":"Relawan Aktif"},{"icon":"🏕️","value":"340+","label":"Kegiatan Terlaksana"},{"icon":"🗺️","value":"28","label":"Kota Terjangkau"},{"icon":"🌳","value":"850T","label":"Pohon Ditanam"}]`},

		// ── Program Unggulan Beranda ─────────────────────────────────────────
		{Key: "home_programs_json", Label: "Daftar Program Unggulan", Section: "Program Unggulan", Type: domain.SettingTypeJSON, OrderNum: 1,
			Value: `[{"icon":"🌊","title":"Bersih Pantai & Sungai","desc":"Aksi bersih rutin di pesisir dan daerah aliran sungai untuk mengurangi sampah plastik di perairan Indonesia."},{"icon":"🌱","title":"Penanaman Mangrove","desc":"Restorasi ekosistem mangrove di wilayah pesisir yang terdegradasi untuk melindungi garis pantai dan habitat satwa."},{"icon":"📚","title":"Edukasi Lingkungan","desc":"Program edukasi lingkungan ke sekolah-sekolah dan komunitas untuk membangun kesadaran sejak dini."}]`},

		// ── Dampak Nyata ─────────────────────────────────────────────────────
		{Key: "home_impacts_json", Label: "Kartu Dampak Nyata", Section: "Dampak Nyata", Type: domain.SettingTypeJSON, OrderNum: 1,
			Value: `[{"icon":"🌳","num":"850.000","unit":"Pohon","desc":"ditanam sejak 2015"},{"icon":"🗑️","num":"120 ton","unit":"Sampah","desc":"berhasil dipungut"},{"icon":"🌿","num":"45 ha","unit":"Mangrove","desc":"berhasil dipulihkan"}]`},

		// ── CTA Bawah Beranda ────────────────────────────────────────────────
		{Key: "home_cta_title", Label: "Judul CTA", Section: "CTA Bawah Beranda", Type: domain.SettingTypeText, OrderNum: 1, Value: "Siap Beraksi untuk Bumi?"},
		{Key: "home_cta_subtitle", Label: "Deskripsi CTA", Section: "CTA Bawah Beranda", Type: domain.SettingTypeText, OrderNum: 2, Value: "Bergabunglah bersama 12.000+ relawan Green Future Indonesia dan jadilah bagian dari perubahan nyata."},

		// ── Halaman Tentang Kami ─────────────────────────────────────────────
		{Key: "about_vision_text", Label: "Teks Visi", Section: "Tentang Kami", Type: domain.SettingTypeText, OrderNum: 1,
			Value: "Mewujudkan Indonesia yang hijau, bersih, dan lestari — di mana manusia dan alam hidup berdampingan secara harmonis untuk generasi yang akan datang."},
		{Key: "about_mission_items", Label: "Poin-Poin Misi (satu poin per baris)", Section: "Tentang Kami", Type: domain.SettingTypeJSON, OrderNum: 2,
			Value: `["Menggerakkan aksi restorasi lingkungan berbasis komunitas","Membangun jaringan relawan lingkungan di seluruh Indonesia","Mendorong kebijakan lingkungan yang berpihak pada alam","Mengedukasi generasi muda tentang krisis iklim","Berkolaborasi dengan pemerintah dan sektor swasta"]`},
		{Key: "about_milestones_json", Label: "Linimasa / Perjalanan Organisasi", Section: "Tentang Kami", Type: domain.SettingTypeJSON, OrderNum: 3,
			Value: `[{"year":"2015","label":"Didirikan di Jakarta oleh 12 aktivis lingkungan muda"},{"year":"2017","label":"Program Penanaman Mangrove perdana di Kepulauan Seribu"},{"year":"2019","label":"Ekspansi ke 10 kota, relawan tembus 3.000 orang"},{"year":"2021","label":"Kemitraan resmi dengan KLHK dan 8 Pemda"},{"year":"2023","label":"Raih penghargaan Lingkungan Nasional dari UNEP Indonesia"},{"year":"2025","label":"12.000+ relawan aktif, 28 kota, 850.000 pohon ditanam"}]`},
		{Key: "about_values_json", Label: "Nilai-Nilai Organisasi", Section: "Tentang Kami", Type: domain.SettingTypeJSON, OrderNum: 4,
			Value: `[{"icon":"🌿","title":"Kepedulian","desc":"Kami percaya setiap tindakan kecil untuk lingkungan memberi dampak besar bagi bumi."},{"icon":"🤝","title":"Kolaborasi","desc":"Perubahan nyata hanya bisa dicapai bersama — lintas komunitas, daerah, dan latar belakang."},{"icon":"🔬","title":"Berbasis Data","desc":"Setiap program dirancang dan dievaluasi berdasarkan data lingkungan yang valid dan terukur."},{"icon":"🌍","title":"Inklusif","desc":"Siapapun bisa berkontribusi — tanpa batasan usia, profesi, atau lokasi."}]`},
		{Key: "about_team_json", Label: "Tim Pengurus", Section: "Tentang Kami", Type: domain.SettingTypeJSON, OrderNum: 5,
			Value: `[{"icon":"👨‍💼","name":"Rizal Anwar, M.Env","role":"Ketua Umum"},{"icon":"👩‍🌾","name":"Nadya Putri, S.Hut","role":"Kepala Program Lapangan"},{"icon":"👨‍💻","name":"Bayu Santoso","role":"Kepala Komunikasi & Media"},{"icon":"👩‍💼","name":"Sari Dewi, M.Si","role":"Manajer Relawan"}]`},

		// ── Informasi Kontak ─────────────────────────────────────────────────
		{Key: "contact_email_1", Label: "Email Utama", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 1, Value: "info@greenfuture.id"},
		{Key: "contact_email_2", Label: "Email Media", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 2, Value: "media@greenfuture.id"},
		{Key: "contact_phone", Label: "Nomor Telepon", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 3, Value: "(021) 456-7890"},
		{Key: "contact_phone_note", Label: "Jam Layanan Telepon", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 4, Value: "Senin–Jumat 08.00–17.00"},
		{Key: "contact_address_1", Label: "Alamat (baris 1)", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 5, Value: "Jl. Kemang Raya No. 45"},
		{Key: "contact_address_2", Label: "Alamat (baris 2 / kota)", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 6, Value: "Jakarta Selatan 12730"},
		{Key: "contact_hours_1", Label: "Jam Operasional (baris 1)", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 7, Value: "Senin – Jumat: 08.00 – 17.00"},
		{Key: "contact_hours_2", Label: "Jam Operasional (baris 2)", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 8, Value: "Sabtu: 09.00 – 13.00 WIB"},
		{Key: "contact_map_label", Label: "Label di Peta", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 9, Value: "Kemang, Jakarta Selatan"},
		{Key: "contact_social_json", Label: "Akun Sosial Media", Section: "Informasi Kontak", Type: domain.SettingTypeJSON, OrderNum: 10,
			Value: `[{"icon":"🌿","platform":"Instagram","handle":"@greenfuture.id"},{"icon":"▶","platform":"YouTube","handle":"Green Future ID"},{"icon":"𝕏","platform":"Twitter/X","handle":"@GreenFutureID"},{"icon":"in","platform":"LinkedIn","handle":"Green Future Indonesia"}]`},

		// ── Footer ───────────────────────────────────────────────────────────
		{Key: "footer_description", Label: "Deskripsi Singkat di Footer", Section: "Footer", Type: domain.SettingTypeText, OrderNum: 1, Value: "Bersama menjaga bumi untuk generasi mendatang melalui aksi nyata, edukasi, dan kolaborasi komunitas."},
		{Key: "footer_programs_json", Label: "Daftar Program di Footer (satu per baris)", Section: "Footer", Type: domain.SettingTypeJSON, OrderNum: 2,
			Value: `["Penanaman Mangrove","Bersih Sungai & Pantai","Edukasi Lingkungan Sekolah","Daur Ulang Sampah Plastik","Pemantauan Kualitas Udara"]`},
		{Key: "footer_social_json", Label: "Ikon Sosial Media di Footer", Section: "Footer", Type: domain.SettingTypeJSON, OrderNum: 3,
			Value: `[{"icon":"𝕏","label":"Twitter/X","link":"#"},{"icon":"in","label":"LinkedIn","link":"#"},{"icon":"f","label":"Facebook","link":"#"},{"icon":"▶","label":"YouTube","link":"#"}]`},

		// ── Halaman Relawan ──────────────────────────────────────────────────
		{Key: "volunteer_why_join_json", Label: "Alasan Bergabung Jadi Relawan", Section: "Halaman Relawan", Type: domain.SettingTypeJSON, OrderNum: 1,
			Value: `[{"icon":"🌿","title":"Dampak Nyata","desc":"Setiap aksi Anda langsung terasa — pohon tertanam, pantai lebih bersih, alam terjaga."},{"icon":"🤝","title":"Komunitas Solid","desc":"Bergabung dengan 12.000+ relawan yang peduli dan saling mendukung."},{"icon":"📜","title":"Sertifikat Resmi","desc":"Dapatkan sertifikat relawan yang diakui dan bisa dicantumkan di CV."},{"icon":"🎓","title":"Pelatihan Gratis","desc":"Akses workshop lingkungan, manajemen komunitas, dan keterampilan lapangan."}]`},
	}

	for _, def := range defs {
		s := def
		db.Where(domain.SiteSetting{Key: s.Key}).FirstOrCreate(&s)
	}
	log.Println("✓ Site settings di-seed")
}
