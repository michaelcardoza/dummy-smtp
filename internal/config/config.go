package config

import (
	"errors"
	"flag"
	"os"

	"github.com/michaelcardoza/dummy-smtp/internal/core/mail"
	"github.com/michaelcardoza/dummy-smtp/internal/infrastructure/storage/memory"
	"github.com/michaelcardoza/dummy-smtp/internal/infrastructure/storage/sqlite"
)

type storageType string

const (
	storageMemory   storageType = "memory"
	storageSQLite   storageType = "sqlite"
	storageMongo    storageType = "mongo"
	storagePostgres storageType = "postgres"
)

type Config struct {
	storageType storageType
	STMPAddr    string
	HTTPAddr    string
	Storage     mail.Storage
	MongoURI    string
	PostgresURI string
}

func defaultConfig() *Config {
	return &Config{
		STMPAddr:    "0.0.0.0:1025",
		HTTPAddr:    "0.0.0.0:8025",
		storageType: storageMemory,
	}
}

func Load() (*Config, error) {
	cfg := defaultConfig()
	registerFlags(cfg)

	switch cfg.storageType {
	case storageMemory:
		cfg.Storage = memory.NewStorage()
	case storageSQLite:
		db, err := sqlite.NewDatabase()
		if err != nil {
			return nil, err
		}

		if err = sqlite.Migrate(db); err != nil {
			return nil, err
		}

		cfg.Storage = sqlite.NewStorage(db)
	default:
		return nil, errors.New("invalid storage")
	}

	return cfg, nil
}

func registerFlags(cfg *Config) {
	flag.StringVar(&cfg.STMPAddr, "smtp-addr", getenv("SMTP_ADDR", "0.0.0.0:1025"), "SMTP port, e.g. 0.0.0.0:1025 or :1025")
	flag.StringVar(&cfg.HTTPAddr, "http-addr", getenv("HTTP_ADDR", "0.0.0.0:8025"), "HTTP port, e.g. 0.0.0.0:8025 or :8025")
	flag.StringVar((*string)(&cfg.storageType), "storage", getenv("STORAGE", "memory"), "Message storage, e.g. memory, sqlite, mongo, postgres")
	flag.StringVar(&cfg.MongoURI, "mongo-uri", getenv("MONGO_URI", "mongodb://127.0.0.1:27017/dummysmtp"), "MongoDB URI, e.g. mongodb://127.0.0.1:27017/dummysmtp")
	flag.StringVar(&cfg.PostgresURI, "pg-uri", getenv("PG_URI", "postgresql://postgres:postgres@127.0.0.1:5432/dummysmtp"), "Postgres URI, e.g. postgresql://postgres:postgres@127.0.0.1:5432/dummysmtp")
	flag.Parse()
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
