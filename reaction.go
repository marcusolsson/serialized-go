package serialized

import (
	"fmt"
	"net/http"
)

// Reaction holds a Serialized.io Reaction.
type Reaction struct {
	ID        string `json:"reactionId"`
	Name      string `json:"reactionName"`
	Feed      string `json:"feedName"`
	EventType string `json:"eventType"`
	Delay     string `json:"delay"`
	Action    Action `json:"action"`
}

// Action holds a Serialized.io Action.
type Action struct {
	HTTPMethod string `json:"httpMethod"`
	TargetURI  string `json:"targetUri"`
	Body       string `json:"body"`
	ActionType string `json:"actionType"`
}

// CreateReaction registers a new reaction.
func (c *Client) CreateReaction(r Reaction) error {
	req, err := c.newRequest("POST", "/reactions", r)
	if err != nil {
		return err
	}

	resp, err := c.do(req, nil)
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return err
}

// ListReactions returns all registered reactions.
func (c *Client) ListReactions() ([]Reaction, error) {
	req, err := c.newRequest("GET", "/reactions", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Reactions []Reaction `json:"reactions"`
	}

	resp, err := c.do(req, &response)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return response.Reactions, err
}

// DeleteReaction deletes a reaction with a given ID.
func (c *Client) DeleteReaction(id string) error {
	req, err := c.newRequest("DELETE", "/reactions/"+id, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Reaction returns a reaction with a given ID.
func (c *Client) Reaction(id string) (Reaction, error) {
	req, err := c.newRequest("DELETE", "/reactions/"+id, nil)
	if err != nil {
		return Reaction{}, err
	}

	var r Reaction
	resp, err := c.do(req, &r)
	if err != nil {
		return Reaction{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Reaction{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return r, nil
}
