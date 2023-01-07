package monzo

// The Monzo app is organised around the feed – a reverse-chronological stream of events.
// Transactions are one such feed item, and your application can create its own feed items to surface relevant information to the user.
type FeedService service

const (
	// FeedTypeBasic is currently the only supported feed type.
	FeedTypeBasic = "basic"
)

// FeedItemParamsBasic represents the customization for the basic feed item type.
type FeedItemParamsBasic struct {
	Title           string `json:"title"`
	ImageURL        string `json:"image_url"`
	Body            string `json:"body,omitempty"`
	BackgroundColor string `json:"background_color,omitempty"`
	TitleColor      string `json:"title_color,omitempty"`
	BodyColor       string `json:"body_color,omitempty"`
}

// FeedItem represents an item that can be posted to a user's feed.
type FeedItem struct {
	AccountID string              `json:"account_id"`
	Type      string              `json:"type"`
	Params    FeedItemParamsBasic `json:"params"`
	URL       string              `json:"url,omitempty"`
}

// Creates a new feed item on the user's feed. These can be dismissed.
func (s *FeedService) Create(feedItem FeedItem) (err error) {
	_, err = s.client.Post("/feed", feedItem)
	return
}
