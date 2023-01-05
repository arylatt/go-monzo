package monzo

// The Monzo app is organised around the feed – a reverse-chronological stream of events.
// Transactions are one such feed item, and your application can create its own feed items to surface relevant information to the user.
type FeedService service
