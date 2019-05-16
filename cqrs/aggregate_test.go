package cqrs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrackChange(t *testing.T) {
	app := New(nil)
	ctx := app.NewContext(context.Background())
	registerTestTypes(app)

	var emits int
	app.Subscribe(func(ctx context.Context, msg Message) error {
		emits++
		return nil
	}, TestEventType)

	event := newTestEvent(ctx,TestEvent{})
	err := app.TrackChange(event)
	if err != nil {
		t.Fatal(err)
	}

	err = app.CommitChanges(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if emits != 1 {
		t.Fatal("expected a single bus emit")
	}

	events := app.Load(ctx,event.Meta().AggregateID, TestAggregateType)

	count := 0
	for range events{
		count++
	}
	if count != 1 {
		t.Fatal("expected a single persisted event")
	}
}

func TestTrackChangeMiddleware(t *testing.T) {
	app := New(nil)
	ctx := app.NewContext(context.Background())
	registerTestTypes(app)

	var emits int
	app.Subscribe(func(ctx context.Context, msg Message) error {
		emits++
		return nil
	}, TestEventType)

	handler := app.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		event := newTestEvent(ctx,TestEvent{&MessageMeta{}})
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
	events := app.Load(ctx,aggregateId, TestAggregateType)

	count := 0
	for range events{
		count++
	}
	if count != 1 {
		t.Fatal("expected a single persisted event")
	}

}

func TestLoadAggregate(t *testing.T) {
	tests := []struct {
		name string
		id string
		err error
	}{
		{
			name: "can load aggregate",
			id: "1234",
			err: nil,
		},
		{
			name: "cant load aggregate if empty id",
			id: "",
			err: ErrNoID,
		},
		{
			name: "cant load aggregate if no events present",
			id: "4321",
			err: ErrNoEvents,
		},
	}

	app := New(nil)
	ctx := app.NewContext(context.Background())
	registerTestTypes(app)

	if err := app.Store(ctx,newTestEvent(ctx,TestEvent{},"1234")); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggr := app.GetAggregate(TestAggregateType)
			err := LoadAggregate(ctx,&AggregateMeta{
				AggregateID:   tt.id,
				AggregateType: TestAggregateType,
			},aggr)

			if err != tt.err {
				t.Fatalf("err was %v, but expected %v",err,tt.err)
			}
		})
	}
}
