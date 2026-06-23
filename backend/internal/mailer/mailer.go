package mailer

import (
	"errors"
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

// Config kredensial SMTP, diisi dari environment variable lewat config.Load()
type Config struct {
	Host     string // misal: smtp.gmail.com
	Port     int    // 587
	Username string // alamat Gmail pengirim
	Password string // App Password Gmail (BUKAN password akun biasa)
	FromName string // nama pengirim yang muncul di email, misal "Yayasan CMS"
}

type Mailer struct {
	cfg Config
}

// ErrSMTPNotConfigured dikembalikan saat SMTP_HOST/USER/PASSWORD belum diisi.
// Dipisah dari error pengiriman biasa supaya handler bisa kasih pesan yang
// berbeda ke admin: "belum dikonfigurasi" vs "gagal kirim" itu dua masalah
// yang beda cara nanganinnya.
var ErrSMTPNotConfigured = errors.New("SMTP belum dikonfigurasi")

func New(cfg Config) *Mailer {
	return &Mailer{cfg: cfg}
}

// IsConfigured mengecek apakah kredensial SMTP sudah diisi.
func (m *Mailer) IsConfigured() bool {
	return m.cfg.Host != "" && m.cfg.Username != "" && m.cfg.Password != ""
}

// SendNewAccountEmail kirim notifikasi akun baru + password default ke user.
//
// BARU: dulu kalau SMTP belum dikonfigurasi, fungsi ini return nil (dianggap
// "berhasil" oleh caller). Sekarang return ErrSMTPNotConfigured supaya
// caller (handler) bisa kasih tahu admin secara jelas alih-alih diam-diam
// menganggap sukses.
func (m *Mailer) SendNewAccountEmail(toEmail, toName, defaultPassword, roleDisplay, loginURL string) error {
	if !m.IsConfigured() {
		log.Printf("⚠ SMTP belum dikonfigurasi — email ke %s TIDAK terkirim (lihat .env: SMTP_HOST/SMTP_USER/SMTP_PASSWORD)", toEmail)
		return ErrSMTPNotConfigured
	}

	subject := "Akun Global Future Indonesia Anda Telah Dibuat"
	body := fmt.Sprintf(`
<div style="font-family:-apple-system,Segoe UI,sans-serif;max-width:480px;margin:0 auto;padding:24px;background:#F8FAFC">
  <div style="background:#fff;border-radius:12px;padding:28px;border:1px solid #E2E8F0">
    <div style="font-size:1.3rem;margin-bottom:4px">🏛️</div>
    <h2 style="color:#0F172A;font-size:1.1rem;margin:0 0 4px">Selamat datang, %s!</h2>
    <p style="color:#64748B;font-size:.85rem;margin:0 0 20px">Akun Anda di Global Future Indonesia telah dibuat oleh administrator dengan role <strong>%s</strong>.</p>

    <div style="background:#F8FAFC;border:1px solid #E2E8F0;border-radius:8px;padding:16px;margin-bottom:20px">
      <table style="width:100%%;font-size:.85rem;color:#334155">
        <tr><td style="padding:4px 0;color:#64748B">Email</td><td style="padding:4px 0;font-weight:600">%s</td></tr>
        <tr><td style="padding:4px 0;color:#64748B">Password Sementara</td><td style="padding:4px 0;font-weight:600;font-family:monospace">%s</td></tr>
      </table>
    </div>

    <div style="background:#FFFBEB;border:1px solid #FDE68A;border-radius:8px;padding:12px 16px;margin-bottom:20px">
      <p style="color:#B45309;font-size:.8rem;margin:0">⚠ <strong>Penting:</strong> Ini adalah password sementara. Demi keamanan, segera login dan ganti password Anda lewat menu <strong>Profil Saya</strong>.</p>
    </div>

    <a href="%s" style="display:inline-block;background:#2563EB;color:#fff;text-decoration:none;padding:10px 20px;border-radius:8px;font-size:.85rem;font-weight:600">Masuk ke CMS</a>

    <p style="color:#94A3B8;font-size:.72rem;margin-top:24px">Email ini dikirim otomatis, mohon tidak membalas ke alamat ini.</p>
  </div>
</div>`, toName, roleDisplay, toEmail, defaultPassword, loginURL)

	return m.send(toEmail, subject, body)
}

func (m *Mailer) send(toEmail, subject, htmlBody string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", m.cfg.FromName, m.cfg.Username))
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", htmlBody)

	dialer := gomail.NewDialer(m.cfg.Host, m.cfg.Port, m.cfg.Username, m.cfg.Password)

	if err := dialer.DialAndSend(msg); err != nil {
		log.Printf("✗ Gagal kirim email ke %s: %v", toEmail, err)
		return fmt.Errorf("gagal kirim email: %w", err)
	}
	log.Printf("✓ Email akun baru terkirim ke %s", toEmail)
	return nil
}

// SendPasswordResetEmail kirim notifikasi password baru hasil reset oleh admin.
//
// BARU: sama seperti SendNewAccountEmail, sekarang return ErrSMTPNotConfigured
// kalau SMTP belum diisi, bukan nil.
func (m *Mailer) SendPasswordResetEmail(toEmail, toName, newPassword, loginURL string) error {
	if !m.IsConfigured() {
		log.Printf("⚠ SMTP belum dikonfigurasi — email reset password ke %s TIDAK terkirim", toEmail)
		return ErrSMTPNotConfigured
	}

	subject := "Password Akun Global Future Indonesia Anda Telah Direset"
	body := fmt.Sprintf(`
<div style="font-family:-apple-system,Segoe UI,sans-serif;max-width:480px;margin:0 auto;padding:24px;background:#F8FAFC">
  <div style="background:#fff;border-radius:12px;padding:28px;border:1px solid #E2E8F0">
    <div style="font-size:1.3rem;margin-bottom:4px">🔑</div>
    <h2 style="color:#0F172A;font-size:1.1rem;margin:0 0 4px">Halo, %s</h2>
    <p style="color:#64748B;font-size:.85rem;margin:0 0 20px">Password akun Anda di Global Future Indonesia telah direset oleh administrator.</p>

    <div style="background:#F8FAFC;border:1px solid #E2E8F0;border-radius:8px;padding:16px;margin-bottom:20px">
      <table style="width:100%%;font-size:.85rem;color:#334155">
        <tr><td style="padding:4px 0;color:#64748B">Password Baru</td><td style="padding:4px 0;font-weight:600;font-family:monospace">%s</td></tr>
      </table>
    </div>

    <div style="background:#FFFBEB;border:1px solid #FDE68A;border-radius:8px;padding:12px 16px;margin-bottom:20px">
      <p style="color:#B45309;font-size:.8rem;margin:0">⚠ Jika Anda tidak meminta reset password ini, segera hubungi administrator.</p>
    </div>

    <a href="%s" style="display:inline-block;background:#2563EB;color:#fff;text-decoration:none;padding:10px 20px;border-radius:8px;font-size:.85rem;font-weight:600">Masuk ke CMS</a>

    <p style="color:#94A3B8;font-size:.72rem;margin-top:24px">Email ini dikirim otomatis, mohon tidak membalas ke alamat ini.</p>
  </div>
</div>`, toName, newPassword, loginURL)

	return m.send(toEmail, subject, body)
}
