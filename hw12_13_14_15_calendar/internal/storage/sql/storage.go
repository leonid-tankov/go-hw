package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/models"
	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/storage/sql/migrations"
	"github.com/pressly/goose/v3"
)

var ErrNoRowsAffected = errors.New("no rows updated")

type Storage struct {
	Dsn string
	DB  *pgx.Conn
}

func New(conf config.Config, logger repository.Logger) *Storage {
	storage := &Storage{
		Dsn: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			conf.Postgres.Username,
			conf.Postgres.Password,
			conf.Postgres.Host,
			conf.Postgres.Port,
			conf.Postgres.Database,
		),
	}
	if err := storage.migrate(logger); err != nil {
		logger.Fatal(err.Error())
	}
	return storage
}

func (s *Storage) migrate(logger repository.Logger) error {
	stdlib.GetDefaultDriver()

	db, err := goose.OpenDBWithDriver("pgx", s.Dsn)
	if err != nil {
		return err
	}
	goose.SetBaseFS(&migrations.EmbedMigrations)

	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	logger.Info("starting migrations...")
	err = goose.Up(db, ".")
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}
	logger.Info("end migrations...")

	return nil
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	s.DB, err = pgx.Connect(ctx, s.Dsn)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.DB.Close(ctx)
}

func (s *Storage) CreateEvent(ctx context.Context, event *models.Event) error {
	statement := `
INSERT INTO events (title, timestamp_event, duration, description, customer_id, notify_duration)
VALUES ($1, $2, $3, $4, $5, $6)
`
	_, err := s.DB.Exec(ctx,
		statement,
		event.Title,
		event.Timestamp.String(),
		event.Duration.Seconds(),
		event.Description,
		event.CustomerID,
		event.NotifyDuration.Seconds())
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	statement := `
SELECT id, title, timestamp_event, duration, description, customer_id, notify_duration FROM events
WHERE id = $1
`
	var event models.Event
	err := s.DB.QueryRow(ctx, statement, id).Scan(
		&event.ID,
		&event.Title,
		&event.Timestamp,
		&event.Duration,
		&event.Description,
		&event.CustomerID,
		&event.NotifyDuration)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *models.Event) error {
	statement := `
UPDATE events
SET title = $2, timestamp_event = $3, duration = $4, description = $5, notify_duration = $6
WHERE id = $1
`
	tag, err := s.DB.Exec(ctx,
		statement,
		event.ID,
		event.Title,
		event.Timestamp.String(),
		event.Duration.Seconds(),
		event.Description,
		event.NotifyDuration.Seconds())
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	statement := `
DELETE FROM events
WHERE id = $1
`
	tag, err := s.DB.Exec(ctx, statement, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (s *Storage) ListEvents(
	ctx context.Context,
	date time.Time,
	days int,
	customerID string,
) ([]*models.Event, error) {
	statement := `
SELECT id, title, timestamp_event, duration, description, customer_id, notify_duration FROM events
WHERE customer_id = $1 and timestamp_event >= $2 and timestamp_event <= $2 + interval '1 day' * $3
`
	rows, err := s.DB.Query(ctx, statement, customerID, date.String(), days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []*models.Event
	for rows.Next() {
		var id, title, description, customer sql.NullString
		var timestamp sql.NullTime
		var duration, notifyDuration sql.NullInt64
		err = rows.Scan(&id, &title, &timestamp, &duration, &description, &customer, &notifyDuration)
		if err != nil {
			return nil, err
		}
		events = append(events, &models.Event{
			ID:             id.String,
			Title:          title.String,
			Timestamp:      timestamp.Time,
			Duration:       time.Duration(duration.Int64) * time.Second,
			Description:    description.String,
			CustomerID:     customer.String,
			NotifyDuration: time.Duration(notifyDuration.Int64) * time.Second,
		})
	}
	return events, nil
}
