package response

// 5.3 Object: Title
//
// Corresponds to the Title Object in the request, with the value filled in.
type Title struct {
	// Field:
	//   text
	// Scope:
	//   required
	// Type:
	//   string
	// Description:
	//   The text associated with the text element.
	Text string `json:"text"`

	// Field:
	//   ext
	// Scope:
	//   optional
	// Type:
	//   object
	// Description:
	//   This object is a placeholder that may contain custom JSON agreed to by the parties to support flexibility beyond the standard defined in this specification
	Ext RawJSON `json:"ext,omitempty"`
}
