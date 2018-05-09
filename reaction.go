package serialized

import (
	"context"
	"fmt"
	"net/http"
)

// Reaction holds a Serialized.io Reaction.
type Reaction struct {
	Name               string   `json:"reactionName,omitempty"`
	Feed               string   `json:"feedName,omitempty"`
	ReactOnEventType   string   `json:"reactOnEventType,omitempty"`
	CancelOnEventTypes []string `json:"cancelOnEventTypes,omitempty"`
	TriggerTimeField   string   `json:"triggerTimeField,omitempty"`
	Offset             string   `json:"offset,omitempty"`
	Action             *Action  `json:"action,omitempty"`
}

// ActionType represents a reaction action.
type ActionType string

// Available action types.
const (
	ActionTypeHTTPPost  ActionType = "HTTP_POST"
	ActionTypeSlackPost ActionType = "SLACK_POST"
)

// Action defines a react action.
type Action struct {
	ActionType ActionType `json:"actionType,omitempty"`
	TargetURI  string     `json:"targetUri,omitempty"`
	Body       string     `json:"body,omitempty"`
}

// CreateReaction registers a new reaction.
func (c *Client) CreateReaction(ctx context.Context, r *Reaction) error {
	req, err := c.newRequest("POST", "/reactions", r)
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, req, nil)
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return err
}

// ListReactions returns all registered reactions.
func (c *Client) ListReactions(ctx context.Context) ([]*Reaction, error) {
	req, err := c.newRequest("GET", "/reactions", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Definitions []*Reaction `json:"definitions"`
	}

	resp, err := c.do(ctx, req, &response)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return response.Definitions, err
}

// DeleteReaction deletes a reaction with a given ID.
func (c *Client) DeleteReaction(ctx context.Context, id string) error {
	req, err := c.newRequest("DELETE", "/reactions/"+id, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Reaction returns a reaction with a given ID.
func (c *Client) Reaction(ctx context.Context, id string) (*Reaction, error) {
	req, err := c.newRequest("DELETE", "/reactions/"+id, nil)
	if err != nil {
		return nil, err
	}

	r := new(Reaction)
	resp, err := c.do(ctx, req, &r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return r, nil
}
