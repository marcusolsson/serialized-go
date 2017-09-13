package serialized

import (
	"fmt"
	"net/http"
)

// Aggregate holds a Serialized.io Aggregate.
type Aggregate struct {
	ID      string  `json:"aggregateId"`
	Version int     `json:"aggregateVersion"`
	Type    string  `json:"aggregateType"`
	Events  []Event `json:"events"`
}

// Store saves events for a given aggregate. All events must refer to
// the same aggregate id.
func (c *Client) Store(aggType, aggID string, version int64, events ...Event) error {
	reqBody := struct {
		AggregateID     string  `json:"aggregateId"`
		Events          []Event `json:"events"`
		ExpectedVersion int64   `json:"expectedVersion,omitempty"`
	}{
		AggregateID:     aggID,
		ExpectedVersion: version,
		Events:          events,
	}

	req, err := c.newRequest("POST", "/aggregates/"+aggType+"/events", reqBody)
	if err != nil {
		return err
	}

	resp, err := c.do(req, nil)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return err
}

// Aggregate exists returns whether a specific aggregate exists.
func (c *Client) AggregateExists(aggType, aggID string) (bool, error) {
	req, err := c.newRequest("HEAD", "/aggregates/"+aggType+"/"+aggID, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.do(req, nil)
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return true, err
}

// LoadAggregate loads all events for a single aggregate.
func (c *Client) LoadAggregate(aggType, aggID string) (Aggregate, error) {
	req, err := c.newRequest("GET", "/aggregates/"+aggType+"/"+aggID, nil)
	if err != nil {
		return Aggregate{}, err
	}

	var a Aggregate
	resp, err := c.do(req, &a)
	if resp.StatusCode != http.StatusOK {
		return Aggregate{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return a, err
}
