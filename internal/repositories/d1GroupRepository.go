package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"

	g "github.com/dannyh79/brp-webhook/internal/groups"
	"github.com/dannyh79/brp-webhook/internal/sentry"
)

var _ Repository[g.Group] = (*D1GroupRepository)(nil)

type endpoint = string

// API token for requesting the endpoint.
type token = string

type SaveGroupParams struct {
	Id string `json:"id"`
}

type D1GroupRepository struct {
	endpoint
	token
	client HttpDoer
}

func (r *D1GroupRepository) Save(g *g.Group) (*g.Group, error) {
	p := SaveGroupParams{Id: g.Id}
	data, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Group: %w", err)
	}

	req, err := http.NewRequest("POST", r.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create save request: %w", err)
	}

	r.setReqHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send save request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

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
	url, err := r.buildReqUrl(g.Id)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create destroy request: %w", err)
	}

	r.setReqHeaders(req)

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send destroy request: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return ErrorNotFound
	default:
		return fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}
}

func (r *D1GroupRepository) setReqHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.token)
}

func (r *D1GroupRepository) buildReqUrl(id string) (string, error) {
	p, err := url.Parse(r.endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}
	p.Path = path.Join(p.Path, id)
	return p.String(), nil
}

func NewD1GroupRepository(u endpoint, t token, c *http.Client) *D1GroupRepository {
	return &D1GroupRepository{endpoint: u, token: t, client: sentry.NewSentryHttpClient(c)}
}
