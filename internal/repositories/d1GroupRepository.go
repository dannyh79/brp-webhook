package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dannyh79/brp-webhook/internal/groups"
)

type endpoint = string

type D1GroupRepository struct {
	endpoint
	client *http.Client
}

func (r *D1GroupRepository) Save(g *groups.Group) (*groups.Group, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Group: %w", err)
	}

	req, err := http.NewRequest("POST", r.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotModified {
		return g, nil
	}

	return nil, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
}

func NewD1GroupRepository(u endpoint, c *http.Client) *D1GroupRepository {
	return &D1GroupRepository{endpoint: u, client: c}
}
