package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// ── Generator password acak yang aman secara kriptografis ───────────────────
//
// PENTING soal keamanan dibanding pendekatan sebelumnya (password default
// statis per role):
//   - Setiap user mendapat password UNIK, bukan dibagi rata per role.
//     Kalau satu password bocor, hanya satu akun yang terdampak.
//   - Memakai crypto/rand (CSPRNG), BUKAN math/rand. math/rand bisa
//     diprediksi kalau seed-nya diketahui/ditebak — tidak aman untuk
//     apapun yang berhubungan dengan kredensial. crypto/rand mengambil
//     entropi dari sistem operasi, tidak bisa ditebak.
//   - Tidak ada nilai password yang disimpan sebagai konstanta di kode
//     sumber — jadi tidak ada apapun untuk "dibaca" dari source/binary.
//   - Password hasil generate HANYA hidup sesaat di memori (return value
//     function ini), langsung di-hash oleh caller, lalu dibuang. Tidak
//     pernah ditulis ke log atau disimpan plaintext di mana pun.

const (
	upperChars   = "ABCDEFGHJKLMNPQRSTUVWXYZ" // tanpa I, O — hindari ambigu dgn 1, 0
	lowerChars   = "abcdefghijkmnopqrstuvwxyz" // tanpa l — hindari ambigu dgn 1
	digitChars   = "23456789"                  // tanpa 0, 1 — hindari ambigu dgn O, l
	symbolChars  = "!@#$%&*?"
	allCharsBase = upperChars + lowerChars + digitChars + symbolChars
)

// GenerateSecurePassword menghasilkan password acak sepanjang `length`
// karakter, dijamin mengandung minimal satu huruf besar, satu huruf kecil,
// satu digit, dan satu simbol — supaya selalu lolos validasi kekuatan
// password standar, sekaligus benar-benar acak.
func GenerateSecurePassword(length int) (string, error) {
	if length < 8 {
		length = 12 // minimum aman, default 12
	}

	// Pastikan minimal 1 karakter dari tiap kategori
	required := []string{
		mustRandomChar(upperChars),
		mustRandomChar(lowerChars),
		mustRandomChar(digitChars),
		mustRandomChar(symbolChars),
	}

	var b strings.Builder
	for _, r := range required {
		b.WriteString(r)
	}

	// Isi sisa panjang dengan karakter acak dari gabungan semua kategori
	for i := len(required); i < length; i++ {
		c, err := randomChar(allCharsBase)
		if err != nil {
			return "", err
		}
		b.WriteString(c)
	}

	// Acak ulang urutan supaya 4 karakter wajib di atas tidak selalu
	// muncul di posisi awal (yang bisa jadi pola yang ditebak)
	shuffled, err := shuffleString(b.String())
	if err != nil {
		return "", err
	}

	return shuffled, nil
}

func randomChar(charset string) (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	if err != nil {
		return "", fmt.Errorf("gagal generate karakter acak: %w", err)
	}
	return string(charset[n.Int64()]), nil
}

// mustRandomChar dipakai untuk kategori wajib saat inisialisasi — kalau
// crypto/rand gagal di sini, ada masalah serius di level OS, jadi panic
// lebih tepat daripada melanjutkan dengan password yang lemah/predictable.
func mustRandomChar(charset string) string {
	c, err := randomChar(charset)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand gagal: %v", err))
	}
	return c
}

func shuffleString(s string) (string, error) {
	runes := []rune(s)
	for i := len(runes) - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		j := n.Int64()
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}
