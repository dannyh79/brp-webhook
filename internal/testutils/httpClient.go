package testutils

import (
	"bytes"
	"io"
	"net/http"
)

type mockHTTPTransport struct {
	statusCode int
}

func (m *mockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString("{}")),
		Header:     make(http.Header),
	}, nil
}

func NewMockHttpClient(mockRespCode int) *http.Client {
	return &http.Client{
		Transport: &mockHTTPTransport{mockRespCode},
	}
}
