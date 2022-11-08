package data

// MatchForDisplay is a struct for collecting and sending data needed for the matches.jet view.
type MatchForDisplay struct {
	OtherUserID   int    `json:"other_user_id"`
	OtherUserName string `json:"other_user_name"`
	MatchID       int    `json:"match_id"`
	PercentMatch  int    `json:"percent_match"`
	ArtistID      int    `json:"artist_id"`
}
