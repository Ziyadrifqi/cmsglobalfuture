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
		if role.ID == 0 {
			log.Fatalf("Role %s gagal tersimpan (ID masih 0)", role.Name)
		}
	}
	log.Println("✓ Role di-seed")
}

func seedAdmin(db *gorm.DB) {
	var role domain.Role
	if err := db.Where("name = ?", domain.RoleSuperAdmin).First(&role).Error; err != nil {
		log.Fatalf("Role super_admin tidak ditemukan saat seed admin: %v", err)
	}

	var existing domain.User
	err := db.Where("email = ?", "admin@yayasan.local").First(&existing).Error
	if err == nil {
		if existing.RoleID != role.ID {
			existing.RoleID = role.ID
			db.Save(&existing)
			log.Println("✓ RoleID admin diperbaiki")
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
		{Slug: "about", Title: "Tentang Kami", Content: "Isi konten tentang kami."},
		{Slug: "vision-mission", Title: "Visi & Misi", Content: "Isi visi dan misi."},
		{Slug: "contact", Title: "Kontak", Content: "Informasi kontak."},
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

// seedSettings mengisi nilai default untuk semua konten yang sebelumnya
// hardcode di kode React (hero, statistik, footer, kontak, tim, dst), supaya
// website tetap tampil sama persis seperti sekarang — tapi sekarang semua
// teks/gambar/list itu bisa diubah dari menu "Konten Website" di CMS.
//
// FirstOrCreate hanya membuat baris baru kalau Key belum ada, jadi aman
// dijalankan berkali-kali tanpa menimpa perubahan yang sudah disimpan admin.
func seedSettings(db *gorm.DB) {
	defs := []domain.SiteSetting{
		// ── Hero - Beranda ──────────────────────────────────────────────────
		{Key: "hero_badge_text", Label: "Teks Badge di Atas Judul", Section: "Hero - Beranda", Type: domain.SettingTypeText, OrderNum: 1,
			Value: "12.000+ Relawan Aktif di 28 Kota"},
		{Key: "hero_title_main", Label: "Judul Hero (baris utama)", Section: "Hero - Beranda", Type: domain.SettingTypeText, OrderNum: 2,
			Value: "Bersama Jaga Bumi"},
		{Key: "hero_title_highlight", Label: "Judul Hero (kata yang disorot kuning)", Section: "Hero - Beranda", Type: domain.SettingTypeText, OrderNum: 3,
			Value: "untuk Generasi Depan"},
		{Key: "hero_subtitle", Label: "Subjudul / Deskripsi Hero", Section: "Hero - Beranda", Type: domain.SettingTypeText, OrderNum: 4,
			Value: "Green Future Indonesia adalah gerakan lingkungan yang mengajak semua orang untuk beraksi nyata — menanam, membersihkan, dan mendidik demi Indonesia yang lebih hijau."},
		{Key: "hero_cta_primary_text", Label: "Teks Tombol Utama", Section: "Hero - Beranda", Type: domain.SettingTypeText, OrderNum: 5,
			Value: "Jadi Relawan"},
		{Key: "hero_cta_secondary_text", Label: "Teks Tombol Kedua", Section: "Hero - Beranda", Type: domain.SettingTypeText, OrderNum: 6,
			Value: "Tentang Kami"},

		// ── Statistik Beranda ───────────────────────────────────────────────
		{Key: "home_stats_json", Label: "Kartu Statistik (value, label, icon)", Section: "Statistik Beranda", Type: domain.SettingTypeJSON, OrderNum: 1,
			Value: `[{"value":"12K+","label":"Relawan Aktif","icon":"🌿"},{"value":"340+","label":"Kegiatan Terlaksana","icon":"🏕️"},{"value":"28","label":"Kota Terjangkau","icon":"🗺️"},{"value":"850T","label":"Pohon Ditanam","icon":"🌳"}]`},

		// ── Program Unggulan ────────────────────────────────────────────────
		{Key: "home_programs_json", Label: "Daftar Program (icon, title, desc)", Section: "Program Unggulan", Type: domain.SettingTypeJSON, OrderNum: 1,
			Value: `[{"icon":"🌊","title":"Bersih Pantai & Sungai","desc":"Aksi bersih rutin di pesisir dan daerah aliran sungai untuk mengurangi sampah plastik di perairan Indonesia."},{"icon":"🌱","title":"Penanaman Mangrove","desc":"Restorasi ekosistem mangrove di wilayah pesisir yang terdegradasi untuk melindungi garis pantai dan habitat satwa."},{"icon":"📚","title":"Edukasi Lingkungan","desc":"Program edukasi lingkungan ke sekolah-sekolah dan komunitas untuk membangun kesadaran sejak dini."}]`},

		// ── Dampak Nyata ────────────────────────────────────────────────────
		{Key: "home_impacts_json", Label: "Kartu Dampak (icon, num, unit, desc)", Section: "Dampak Nyata", Type: domain.SettingTypeJSON, OrderNum: 1,
			Value: `[{"icon":"🌳","num":"850.000","unit":"Pohon","desc":"ditanam sejak 2015"},{"icon":"🗑️","num":"120 ton","unit":"Sampah","desc":"berhasil dipungut"},{"icon":"🌿","num":"45 ha","unit":"Mangrove","desc":"berhasil dipulihkan"}]`},

		// ── CTA Bawah Beranda ───────────────────────────────────────────────
		{Key: "home_cta_title", Label: "Judul", Section: "CTA Bawah Beranda", Type: domain.SettingTypeText, OrderNum: 1,
			Value: "Siap Beraksi untuk Bumi?"},
		{Key: "home_cta_subtitle", Label: "Deskripsi", Section: "CTA Bawah Beranda", Type: domain.SettingTypeText, OrderNum: 2,
			Value: "Bergabunglah bersama 12.000+ relawan Green Future Indonesia dan jadilah bagian dari perubahan nyata."},

		// ── Footer ──────────────────────────────────────────────────────────
		{Key: "footer_description", Label: "Deskripsi Singkat", Section: "Footer", Type: domain.SettingTypeText, OrderNum: 1,
			Value: "Bersama menjaga bumi untuk generasi mendatang melalui aksi nyata, edukasi, dan kolaborasi komunitas."},
		{Key: "footer_programs_json", Label: "Daftar Nama Program (array teks)", Section: "Footer", Type: domain.SettingTypeJSON, OrderNum: 2,
			Value: `["Penanaman Mangrove","Bersih Sungai & Pantai","Edukasi Lingkungan Sekolah","Daur Ulang Sampah Plastik","Pemantauan Kualitas Udara"]`},
		{Key: "footer_social_json", Label: "Sosial Media (icon, label, link)", Section: "Footer", Type: domain.SettingTypeJSON, OrderNum: 3,
			Value: `[{"icon":"𝕏","label":"Twitter/X","link":"#"},{"icon":"in","label":"LinkedIn","link":"#"},{"icon":"f","label":"Facebook","link":"#"},{"icon":"▶","label":"YouTube","link":"#"}]`},

		// ── Informasi Kontak (dipakai Footer & Halaman Kontak) ─────────────
		{Key: "contact_email_1", Label: "Email Utama", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 1,
			Value: "info@greenfuture.id"},
		{Key: "contact_email_2", Label: "Email Media", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 2,
			Value: "media@greenfuture.id"},
		{Key: "contact_phone", Label: "Nomor Telepon", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 3,
			Value: "(021) 456-7890"},
		{Key: "contact_phone_note", Label: "Catatan Jam Telepon", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 4,
			Value: "Senin–Jumat 08.00–17.00"},
		{Key: "contact_address_1", Label: "Alamat (baris 1)", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 5,
			Value: "Jl. Kemang Raya No. 45"},
		{Key: "contact_address_2", Label: "Alamat (baris 2)", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 6,
			Value: "Jakarta Selatan 12730"},
		{Key: "contact_hours_1", Label: "Jam Operasional (baris 1)", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 7,
			Value: "Senin – Jumat: 08.00 – 17.00"},
		{Key: "contact_hours_2", Label: "Jam Operasional (baris 2)", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 8,
			Value: "Sabtu: 09.00 – 13.00 WIB"},
		{Key: "contact_map_label", Label: "Label di Bawah Peta", Section: "Informasi Kontak", Type: domain.SettingTypeText, OrderNum: 9,
			Value: "Kemang, Jakarta Selatan"},
		{Key: "contact_social_json", Label: "Sosial Media Halaman Kontak (platform, handle, icon)", Section: "Informasi Kontak", Type: domain.SettingTypeJSON, OrderNum: 10,
			Value: `[{"platform":"Instagram","handle":"@greenfuture.id","icon":"🌿"},{"platform":"YouTube","handle":"Green Future ID","icon":"▶"},{"platform":"Twitter/X","handle":"@GreenFutureID","icon":"𝕏"},{"platform":"LinkedIn","handle":"Green Future Indonesia","icon":"in"}]`},

		// ── Halaman Tentang Kami ────────────────────────────────────────────
		{Key: "about_team_json", Label: "Tim Pengurus (name, role, icon)", Section: "Halaman Tentang Kami", Type: domain.SettingTypeJSON, OrderNum: 1,
			Value: `[{"name":"Rizal Anwar, M.Env","role":"Ketua Umum","icon":"👨‍💼"},{"name":"Nadya Putri, S.Hut","role":"Kepala Program Lapangan","icon":"👩‍🌾"},{"name":"Bayu Santoso","role":"Kepala Komunikasi & Media","icon":"👨‍💻"},{"name":"Sari Dewi, M.Si","role":"Manajer Relawan","icon":"👩‍💼"}]`},
		{Key: "about_values_json", Label: "Nilai-Nilai (icon, title, desc)", Section: "Halaman Tentang Kami", Type: domain.SettingTypeJSON, OrderNum: 2,
			Value: `[{"icon":"🌿","title":"Kepedulian","desc":"Kami percaya setiap tindakan kecil untuk lingkungan memberi dampak besar bagi bumi."},{"icon":"🤝","title":"Kolaborasi","desc":"Perubahan nyata hanya bisa dicapai bersama — lintas komunitas, daerah, dan latar belakang."},{"icon":"🔬","title":"Berbasis Data","desc":"Setiap program dirancang dan dievaluasi berdasarkan data lingkungan yang valid dan terukur."},{"icon":"🌍","title":"Inklusif","desc":"Siapapun bisa berkontribusi — tanpa batasan usia, profesi, atau lokasi."}]`},
		{Key: "about_milestones_json", Label: "Linimasa / Milestone (year, label)", Section: "Halaman Tentang Kami", Type: domain.SettingTypeJSON, OrderNum: 3,
			Value: `[{"year":"2015","label":"Didirikan di Jakarta oleh 12 aktivis lingkungan muda"},{"year":"2017","label":"Program Penanaman Mangrove perdana di Kepulauan Seribu"},{"year":"2019","label":"Ekspansi ke 10 kota, relawan tembus 3.000 orang"},{"year":"2021","label":"Kemitraan resmi dengan KLHK dan 8 Pemda"},{"year":"2023","label":"Raih penghargaan Lingkungan Nasional dari UNEP Indonesia"},{"year":"2025","label":"12.000+ relawan aktif, 28 kota, 850.000 pohon ditanam"}]`},

		// ── Halaman Relawan ─────────────────────────────────────────────────
		{Key: "volunteer_why_join_json", Label: "Kenapa Bergabung (icon, title, desc)", Section: "Halaman Relawan", Type: domain.SettingTypeJSON, OrderNum: 1,
			Value: `[{"icon":"🌿","title":"Dampak Nyata","desc":"Setiap aksi Anda langsung terasa — pohon tertanam, pantai lebih bersih, alam terjaga."},{"icon":"🤝","title":"Komunitas Solid","desc":"Bergabung dengan 12.000+ relawan yang peduli dan saling mendukung."},{"icon":"📜","title":"Sertifikat Resmi","desc":"Dapatkan sertifikat relawan yang diakui dan bisa dicantumkan di CV."},{"icon":"🎓","title":"Pelatihan Gratis","desc":"Akses workshop lingkungan, manajemen komunitas, dan keterampilan lapangan."}]`},

		// ── Identitas Situs (logo & nama, dipakai Navbar) ──────────────────
		{Key: "site_name", Label: "Nama Situs (baris atas)", Section: "Identitas Situs", Type: domain.SettingTypeText, OrderNum: 1,
			Value: "Green Future"},
		{Key: "site_name_sub", Label: "Nama Situs (baris bawah, kecil)", Section: "Identitas Situs", Type: domain.SettingTypeText, OrderNum: 2,
			Value: "Indonesia"},
		{Key: "site_logo_image", Label: "Logo (opsional — kosongkan untuk pakai ikon 🌿 bawaan)", Section: "Identitas Situs", Type: domain.SettingTypeImage, OrderNum: 3,
			Value: ""},
	}

	for _, def := range defs {
		s := def
		db.Where(domain.SiteSetting{Key: s.Key}).FirstOrCreate(&s)
	}
	log.Println("✓ Site settings (konten dinamis) di-seed")
}
