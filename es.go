package serialized

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// ErrAggregateNotFound is returned when no events exist for a given aggregate ID.
var ErrAggregateNotFound = errors.New("aggregate not found")

// Aggregate holds a Serialized.io Aggregate.
type Aggregate struct {
	ID      string   `json:"aggregateId"`
	Version int      `json:"aggregateVersion"`
	Type    string   `json:"aggregateType"`
	Events  []*Event `json:"events"`
}

// Store saves events for a given aggregate. All events must refer to
// the same aggregate id.
func (c *Client) Store(ctx context.Context, aggType, aggID string, version int64, events ...*Event) error {
	reqBody := struct {
		AggregateID     string   `json:"aggregateId"`
		Events          []*Event `json:"events"`
		ExpectedVersion int64    `json:"expectedVersion,omitempty"`
	}{
		AggregateID:     aggID,
		ExpectedVersion: version,
		Events:          events,
	}

	req, err := c.newRequest("POST", "/aggregates/"+aggType+"/events", reqBody)
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

// RequestDeleteAggregateByType requests a aggregate deletion. To delete the
// aggregates, pass the token returned by this this method to
// DeleteAggregateByType.
func (c *Client) RequestDeleteAggregateByType(ctx context.Context, aggType string) (string, error) {
	req, err := c.newRequest("DELETE", "/aggregates/"+aggType, nil)
	if err != nil {
		return "", err
	}

	var response struct {
		DeleteToken string `json:"deleteToken"`
	}

	resp, err := c.do(ctx, req, &response)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return response.DeleteToken, nil
}

// DeleteAggregateByType permanently deletes all aggregates of a given type.
// It requires a token returned from RequestDeleteAggregateByType.
func (c *Client) DeleteAggregateByType(ctx context.Context, aggType, token string) error {
	path := fmt.Sprintf("/aggregates/%s?deleteToken=%s", aggType, token)

	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// AggregateExists returns whether a specific aggregate exists.
func (c *Client) AggregateExists(ctx context.Context, aggType, aggID string) (bool, error) {
	req, err := c.newRequest("HEAD", "/aggregates/"+aggType+"/"+aggID, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.do(ctx, req, nil)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return true, nil
}

// LoadAggregate loads all events for a single aggregate.
func (c *Client) LoadAggregate(ctx context.Context, aggType, aggID string) (*Aggregate, error) {
	req, err := c.newRequest("GET", "/aggregates/"+aggType+"/"+aggID, nil)
	if err != nil {
		return nil, err
	}

	a := new(Aggregate)
	resp, err := c.do(ctx, req, a)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return a, nil
}
