package serialized

import (
	"context"
	"fmt"
	"net/http"
)

// ReactionDefinition holds a Serialized.io Reaction.
type ReactionDefinition struct {
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

// CreateReactionDefinition registers a new reaction definition.
func (c *Client) CreateReactionDefinition(ctx context.Context, r *ReactionDefinition) error {
	req, err := c.newRequest("POST", "/reactions/definitions", r)
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, req, nil)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return err
}

// ListReactionDefinitions returns all registered reactions.
func (c *Client) ListReactionDefinitions(ctx context.Context) ([]*ReactionDefinition, error) {
	req, err := c.newRequest("GET", "/reactions/definitions", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Definitions []*ReactionDefinition `json:"definitions"`
	}

	resp, err := c.do(ctx, req, &response)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return response.Definitions, err
}

// DeleteReactionDefinition deletes a reaction with a given name.
func (c *Client) DeleteReactionDefinition(ctx context.Context, name string) error {
	req, err := c.newRequest("DELETE", "/reactions/definitions/"+name, nil)
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

// ReactionDefinition returns a reaction definition with a given name.
func (c *Client) ReactionDefinition(ctx context.Context, name string) (*ReactionDefinition, error) {
	req, err := c.newRequest("GET", "/reactions/definitions/"+name, nil)
	if err != nil {
		return nil, err
	}

	r := new(ReactionDefinition)
	resp, err := c.do(ctx, req, &r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return r, nil
}
