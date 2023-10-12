package repository

import (
	"context"

	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/models"
)

type Application interface {
	CreateEvent(ctx context.Context, event *models.Event) error
}
