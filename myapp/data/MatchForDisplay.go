// Spotmeet - (Capstone Team E)
// 2022 Stephen Sumpter, John Neumeier,
// Zach Kods, Landon Wilson
package data

// MatchForDisplay is a struct for collecting and sending data needed for the matches.jet view.
type MatchForDisplay struct {
	OtherUserID   int    `json:"other_user_id"`
	OtherUserName string `json:"other_user_name"`
	MatchID       int    `json:"match_id"`
	PercentMatch  int    `json:"percent_match"`
	SongID        int    `json:"song_id"`
	SongName      string `json:"song_name"`
	ArtistName    string `json:"artist_name"`
}
