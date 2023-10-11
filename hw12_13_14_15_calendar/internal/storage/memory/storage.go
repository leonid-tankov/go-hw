package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/models"
)

var (
	ErrEventExist    = errors.New("there is already such an event")
	ErrEventNotExist = errors.New("there is no such an event")
	ErrPastTime      = errors.New("the event is in the past")
)

type memoryStorage map[string]*models.Event

type Storage struct {
	mu   sync.RWMutex
	data memoryStorage
}

func New() *Storage {
	return &Storage{
		mu:   sync.RWMutex{},
		data: make(memoryStorage),
	}
}

func (s *Storage) CreateEvent(ctx context.Context, event *models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[event.ID]; ok {
		return ErrEventExist
	}
	if event.Timestamp.Before(time.Now()) {
		return ErrPastTime
	}
	s.data[event.ID] = event
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[id]; !ok {
		return nil, nil
	}
	event := s.data[id]
	return event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[event.ID]; !ok {
		return ErrEventNotExist
	}
	if event.Timestamp.Before(time.Now()) {
		return ErrPastTime
	}
	s.data[event.ID] = event
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[id]; !ok {
		return ErrEventNotExist
	}
	delete(s.data, id)
	return nil
}

func (s *Storage) ListEvents(
	ctx context.Context,
	date time.Time,
	days int,
	customerID string,
) ([]*models.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var events []*models.Event
	for _, event := range s.data {
		if event.CustomerID == customerID && event.Timestamp.After(date) && event.Timestamp.Before(date.AddDate(0, 0, days)) {
			events = append(events, event)
		}
	}
	return events, nil
}
