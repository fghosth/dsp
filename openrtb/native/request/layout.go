package request

// 7.1 Native Layout IDs - To Be Deprecated
//
// Layout ID is to be deprecated in a future version and is not suggested for new implementations.
//
// Below is a list of the core layouts described in the introduction above.
//
// An implementing exchange may not support all asset variants or introduce new ones unique to that system.
type Layout int64

const (
	LayoutContentWall   Layout = 1 // Content Wall
	LayoutAppWall       Layout = 2 // App Wall
	LayoutNewsFeed      Layout = 3 // News Feed
	LayoutChatList      Layout = 4 // Chat List
	LayoutCarousel      Layout = 5 // Carousel
	LayoutContentStream Layout = 6 // Content Stream
	LayoutGrid          Layout = 7 // Grid adjoining the content

// 500+ Reserved for Exchange specific layouts.
)
