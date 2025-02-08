package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	g "github.com/dannyh79/brp-webhook/internal/groups"
)

var _ Repository[g.Group] = (*D1GroupRepository)(nil)

type endpoint = string

type SaveGroupParams struct {
	Id string `json:"id"`
}

type D1GroupRepository struct {
	endpoint
	client *http.Client
}

func (r *D1GroupRepository) Save(g *g.Group) (*g.Group, error) {
	p := SaveGroupParams{Id: g.Id}
	data, err := json.Marshal(p)
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

	switch resp.StatusCode {
	case http.StatusNoContent:
		return g, nil
	case http.StatusNotModified:
		return nil, ErrorAlreadyExists
	default:
		return nil, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}
}

func (r *D1GroupRepository) Destroy(g *g.Group) error {
	return nil
}

func NewD1GroupRepository(u endpoint, c *http.Client) *D1GroupRepository {
	return &D1GroupRepository{endpoint: u, client: c}
}
