package cqrs

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewBufferedResponseWriter(t *testing.T) {
	recorder := httptest.NewRecorder()
	buf := NewBufferedResponseWriter(recorder)
	buf.Header().Set("X-FOO","/bar")
	buf.Write([]byte("Hello world!"))
	buf.WriteHeader(http.StatusTemporaryRedirect)

	if err := buf.Close(); err != nil {
		log.Fatal(err)
	}

	result := recorder.Result()
	if result.Header.Get("X-FOO") != "/bar"{
		t.Fatal("expected header to be present")
	}

	if recorder.Body.String() != "Hello world!" {
		t.Fatal("expected body to be present")
	}

}