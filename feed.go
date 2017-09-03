package serialized

// Feed holds a Serialized.io feed.
type Feed struct {
	Entries []FeedEntry `json:"entries"`
	HasMore bool        `json:"hasMore"`
}

// FeedEntry holds a Serialized.io feed entry.
type FeedEntry struct {
	SequenceNumber int64
	AggregateID    string
	Timestamp      int64
	Events         []Event
}

// Feeds returns all feed types.
func (c *Client) Feeds() ([]string, error) {
	req, err := c.newRequest("GET", "/feeds/", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Feeds []string `json:"feeds"`
	}

	_, err = c.do(req, &response)
	return response.Feeds, err
}

// Feed returns the feed for a given aggregate.
func (c *Client) Feed(name string) (Feed, error) {
	req, err := c.newRequest("GET", "/feeds/"+name, nil)
	if err != nil {
		return Feed{}, err
	}

	var f Feed
	_, err = c.do(req, &f)
	return f, err
}
