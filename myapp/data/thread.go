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
