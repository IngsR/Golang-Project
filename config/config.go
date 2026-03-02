package config

// AppConfig menyimpan seluruh konfigurasi aplikasi
type AppConfig struct {
	Port       string // Port server HTTP
	DBPath     string // Path file database SQLite
	SecretKey  string // Secret key untuk session
	AppName    string // Nama aplikasi
}

// LoadConfig mengembalikan konfigurasi default aplikasi
func LoadConfig() *AppConfig {
	return &AppConfig{
		Port:      ":8080",
		DBPath:    "goproject.db",
		SecretKey: "rahasia-session-key-ganti-di-production",
		AppName:   "Sistem Informasi Go",
	}
}
