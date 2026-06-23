package domain

import "time"

// SettingType menentukan bagaimana value sebuah setting harus diperlakukan
// di form CMS maupun saat dipakai di frontend.
type SettingType string

const (
	SettingTypeText  SettingType = "text"  // teks pendek/panjang biasa
	SettingTypeImage SettingType = "image" // path gambar hasil upload
	SettingTypeJSON  SettingType = "json"  // list/array, mis. statistik, anggota tim, dll
)

// SiteSetting menyimpan konten dinamis portal (teks hero, statistik, info
// kontak, footer, tim pengurus, dst) supaya bisa diubah dari CMS tanpa
// perlu deploy ulang frontend. Setiap baris punya Key unik yang dikenali
// oleh kode frontend, dikelompokkan per Section supaya enak ditampilkan
// di form CMS.
type SiteSetting struct {
	ID        uint        `gorm:"primaryKey"`
	Key       string      `gorm:"uniqueIndex;not null"` // identifier unik, dipakai frontend, JANGAN diubah-ubah
	Label     string      `gorm:"not null"`             // label yang tampil di form CMS
	Section   string      `gorm:"not null;index"`       // pengelompokan section di form CMS, mis. "Hero - Beranda"
	Type      SettingType `gorm:"default:'text';column:value_type"`
	Value     string      `gorm:"type:text"`
	OrderNum  int         `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
