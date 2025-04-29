package config

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DB        *gorm.DB
	JWTSecret string
}

func LoadConfig() *Config {
	godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(gormpg.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrasyonları çalıştır
	err = runMigrations(db, dsn)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	return &Config{
		DB:        db,
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}

func runMigrations(db *gorm.DB, dsn string) error {
	// GORM'un *gorm.DB'sinden *sql.DB al
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Migrate için PostgreSQL sürücüsü
	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		return err
	}

	// Migrasyon nesnesi oluştur
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Migrasyon dosyalarının yolu
		"postgres",          // Veritabanı adı
		driver,
	)
	if err != nil {
		return err
	}

	// Migrasyonları uygula
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
