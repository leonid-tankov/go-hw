package storage

import (
	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/repository"
	memorystorage "github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var (
	MemoryStorage   = "memory"
	PostgresStorage = "postgres"
)

func NewStorageByType(conf config.Config, logger repository.Logger) repository.Storage {
	switch conf.Storage.Type {
	case MemoryStorage:
		return memorystorage.New()
	case PostgresStorage:
		return sqlstorage.New(conf, logger)
	default:
		return memorystorage.New()
	}
}
