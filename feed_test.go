package monzo

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFeedCreate(t *testing.T) {
	expected := FeedItem{
		AccountID: "test",
		Type:      FeedTypeBasic,
		URL:       "https://www.google.com",
		Params: FeedItemParamsBasic{
			Title:    "Hello, world!",
			ImageURL: "https://www.nyan.cat/cats/original.gif",
		},
	}

	c := MockRequest(nil, func(args mock.Arguments) {
		req := args.Get(0).(*http.Request)

		assert.Equal(t, "/feed", req.URL.Path)

		actual := FeedItem{}

		assert.NoError(t, json.NewDecoder(req.Body).Decode(&actual))
		assert.Equal(t, expected, actual)
	})

	assert.NoError(t, c.Feed.Create(expected))
}
