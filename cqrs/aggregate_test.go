package cqrs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrackChange(t *testing.T) {
	store := NewMemoryEventStore()
	bus := NewLocalMessageBus()
	ct := NewChangeTracker(store, bus)

	var emits int
	bus.Subscribe(func(ctx context.Context, msg Message) error {
		emits++
		return nil
	}, TestEventType)

	event := NewTestEvent(TestEvent{})
	err := ct.TrackChange(event)
	if err != nil {
		t.Fatal(err)
	}

	err = ct.CommitChanges(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if emits != 1 {
		t.Fatal("expected a single bus emit")
	}

	events, err := store.Load(event.Meta().AggregateID, TestAggregateType)
	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 1 {
		t.Fatal("expected a single persisted event")
	}
}

func TestTrackChangeMiddleware(t *testing.T) {
	store := NewMemoryEventStore()
	bus := NewLocalMessageBus()
	ct := NewChangeTracker(store, bus)

	var emits int
	bus.Subscribe(func(ctx context.Context, msg Message) error {
		emits++
		return nil
	}, TestEventType)

	handler := ct.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		event := NewTestEvent(TestEvent{})
		err := DispatchMessage(r.Context(), event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(event.Meta().AggregateID))
	}))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/", nil)
	handler.ServeHTTP(recorder, request)

	if emits != 1 {
		t.Fatal("expected a single bus emit")
	}

	aggregateId := recorder.Body.String()
	events, err := store.Load(aggregateId, TestAggregateType)
	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 1 {
		t.Fatal("expected a single persisted event")
	}

}
