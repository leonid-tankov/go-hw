package repository

import (
	"context"
	"time"

	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/models"
)

type Storage interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEvent(ctx context.Context, id string) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, id string) error
	ListEvents(ctx context.Context, date time.Time, days int, customerID string) ([]*models.Event, error)
}
