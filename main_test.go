package main

import (
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {

	req := httptest.NewRequest("GET", "/user/john@example.com/id-123", nil)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	rec := httptest.NewRecorder()

	r := makeRouter()
	r.ServeHTTP(rec, req)

	if email := rec.Body.String(); email != "john@example.com" {
		t.Errorf("router failed, expected: %s, got: %s", "john@example.com", email)
	}
}
