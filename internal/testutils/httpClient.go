package testutils

import (
	"bytes"
	"io"
	"net/http"
)

type Callback = func(req *http.Request)

type mockHTTPTransport struct {
	statusCode int
	callbacks  []Callback
}

func (m *mockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, cb := range m.callbacks {
		cb(req)
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString("{}")),
		Header:     make(http.Header),
	}, nil
}

func NewMockHttpClient(mockRespCode int, callbacks ...Callback) *http.Client {
	return &http.Client{
		Transport: &mockHTTPTransport{mockRespCode, callbacks},
	}
}
