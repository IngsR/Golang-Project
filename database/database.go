package database

import (
	"goproject/models"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// DB adalah variabel global untuk koneksi database
var DB *gorm.DB

// InitDB menginisialisasi koneksi database SQLite dan menjalankan auto-migration
func InitDB(dbPath string) *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	// Auto-migrate model ke tabel database
	err = DB.AutoMigrate(&models.User{}, &models.Article{})
	if err != nil {
		log.Fatal("Gagal migrasi database:", err)
	}

	log.Println("Database berhasil terkoneksi dan migrasi selesai")
	return DB
}
