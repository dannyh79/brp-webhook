package testutils

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertHttpStatus(t *testing.T) func(rr *httptest.ResponseRecorder, want int) {
	return func(rr *httptest.ResponseRecorder, want int) {
		t.Helper()
		got := rr.Result().StatusCode
		assert.Equal(t, want, got, "got HTTP status %v, want %v", got, want)
	}
}
