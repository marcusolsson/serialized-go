package serialized

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Projection struct {
	ID   string          `json:"projectionId"`
	Data json.RawMessage `json:"data"`
}

type ProjectionDefinition struct {
	Name     string          `json:"projectionName"`
	Feed     string          `json:"feedName"`
	Handlers []*EventHandler `json:"handlers"`
}

type EventHandler struct {
	EventType string      `json:"eventType"`
	Functions []*Function `json:"functions"`
}

type Function struct {
	Name           string `json:"name"`
	TargetSelector string `json:"targetSelector"`
	EventSelector  string `json:"eventSelector"`
	TargetFilter   string `json:"targetFilter"`
	EventFilter    string `json:"eventFilter"`
}

func (c *Client) ListProjectionDefinitions(ctx context.Context) ([]*ProjectionDefinition, error) {
	req, err := c.newRequest("GET", "/projections/definitions", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Definitions []*ProjectionDefinition `json:"definitions"`
	}

	resp, err := c.do(ctx, req, &response)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return response.Definitions, err
}

func (c *Client) CreateProjectionDefinition(ctx context.Context, d *ProjectionDefinition) error {
	req, err := c.newRequest("POST", "/projections/definitions", d)
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, req, nil)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return err
}

func (c *Client) Projection(ctx context.Context, projName, aggID string) (*Projection, error) {
	req, err := c.newRequest("GET", "/projections/single/"+projName+"/"+aggID, nil)
	if err != nil {
		return nil, err
	}

	var proj Projection

	resp, err := c.do(ctx, req, &proj)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return &proj, err
}
