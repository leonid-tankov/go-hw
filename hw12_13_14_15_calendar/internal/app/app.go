package app

import (
	"context"

	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/models"
	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/repository"
)

type App struct {
	logger  repository.Logger
	storage repository.Storage
}

func New(logger repository.Logger, storage repository.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event *models.Event) error {
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}
