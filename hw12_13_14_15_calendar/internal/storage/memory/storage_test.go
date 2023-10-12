package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/leonid-tankov/go-hw/hw12_13_14_15_calendar/internal/models"
	"github.com/stretchr/testify/require"
)

var (
	firstEvent = &models.Event{
		ID:             "1",
		Title:          "title1",
		Duration:       time.Second,
		Description:    "description1",
		CustomerID:     "1",
		NotifyDuration: time.Second,
	}
	secondEvent = &models.Event{
		ID:             "2",
		Title:          "title2",
		Duration:       time.Second,
		Description:    "description2",
		CustomerID:     "2",
		NotifyDuration: time.Second,
	}
	thirdEvent = &models.Event{
		ID:             "3",
		Title:          "title3",
		Duration:       time.Second,
		Description:    "description3",
		CustomerID:     "2",
		NotifyDuration: time.Second,
	}
)

func TestStorage(t *testing.T) {
	store := New()

	event, err := store.GetEvent(context.Background(), "1")
	require.NoError(t, err)
	require.Nil(t, event)

	events, err := store.ListEvents(context.Background(), time.Now(), 1, "1")
	require.NoError(t, err)
	require.Nil(t, events)

	firstEvent.Timestamp = time.Now().Add(-time.Second)
	err = store.CreateEvent(context.Background(), firstEvent)
	require.ErrorIs(t, err, ErrPastTime)

	firstEvent.Timestamp = time.Now().Add(time.Second * 10)
	err = store.CreateEvent(context.Background(), firstEvent)
	require.NoError(t, err)
	err = store.CreateEvent(context.Background(), firstEvent)
	require.ErrorIs(t, err, ErrEventExist)

	event, err = store.GetEvent(context.Background(), firstEvent.ID)
	require.NoError(t, err)
	require.Equal(t, firstEvent, event)

	err = store.UpdateEvent(context.Background(), secondEvent)
	require.ErrorIs(t, err, ErrEventNotExist)

	firstEvent.Timestamp = time.Now().Add(-time.Second)
	err = store.UpdateEvent(context.Background(), firstEvent)
	require.ErrorIs(t, err, ErrPastTime)

	err = store.DeleteEvent(context.Background(), secondEvent.ID)
	require.ErrorIs(t, err, ErrEventNotExist)

	secondEvent.Timestamp = time.Now().Add(time.Hour)
	err = store.CreateEvent(context.Background(), secondEvent)
	require.NoError(t, err)
	err = store.DeleteEvent(context.Background(), secondEvent.ID)
	require.NoError(t, err)

	secondEvent.Timestamp = time.Now().Add(time.Hour * 48) // 2 days
	err = store.CreateEvent(context.Background(), secondEvent)
	require.NoError(t, err)
	thirdEvent.Timestamp = time.Now().Add(time.Hour * 96) // 4 days
	secondEvent.ID = "3"
	err = store.CreateEvent(context.Background(), thirdEvent)
	require.NoError(t, err)
	events, err = store.ListEvents(context.Background(), time.Now(), 1, "2")
	require.NoError(t, err)
	require.Equal(t, 0, len(events))
	events, err = store.ListEvents(context.Background(), time.Now(), 3, "2")
	require.NoError(t, err)
	require.Equal(t, 1, len(events))
	events, err = store.ListEvents(context.Background(), time.Now(), 5, "2")
	require.NoError(t, err)
	require.Equal(t, 2, len(events))
}
