package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_APIServer_handleHello(t *testing.T) {
	s := New(NewConfig())
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/hello", nil)
	s.handleHello().ServeHTTP(rec, req)
	if rec.Body.String() != "Hello" {
		t.Errorf("handleHello -> expected = %s, actual = %s", "Hello", rec.Body.String())
	}
}
