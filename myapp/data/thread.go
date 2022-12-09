// Spotmeet - (Capstone Team E)
// 2022 Stephen Sumpter, John Neumeier,
// Zach Kods, Landon Wilson
package data

import "time"

type Thread struct {
	UserID                int
	MatchID               int
	MatchFirstName        string
	LatestMessagePreview  string
	LatestMessageTimeSent string
	OtherUsersImage       string
	TimeSentISO           time.Time
	UserHasOpenedThread   bool
}
