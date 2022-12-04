// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import "testing"

func TestMessage_Insert(t *testing.T) {
	u1 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "message1@test.com",
		Active:    1,
	}

	u1ID, err := u1.Insert(u1)
	if err != nil {
		t.Error(err)
	}

	u2 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "message2@test.com",
		Active:    1,
	}

	u2ID, err := u2.Insert(u2)
	if err != nil {
		t.Error(err)
	}

	m := Message{
		UserID:  u1ID,
		MatchID: u2ID,
		Content: "test",
	}

	mID, err := m.Insert(m)
	if err != nil {
		t.Error(err)
	}

	if mID == 0 {
		t.Error("failed to insert message")
	}

}

func TestMessage_Get(t *testing.T) {
	u1 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "message3@test.com",
		Active:    1,
	}

	u1ID, err := u1.Insert(u1)
	if err != nil {
		t.Error(err)
	}

	u2 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "message4@test.com",
		Active:    1,
	}

	u2ID, err := u2.Insert(u2)
	if err != nil {
		t.Error(err)
	}

	m := Message{
		UserID:  u1ID,
		MatchID: u2ID,
		Content: "test",
	}

	mID, err := m.Insert(m)
	if err != nil {
		t.Error(err)
	}

	m.ID = mID

	message, err := m.Get(mID)
	if err != nil {
		t.Error(err)
	}

	if message.ID != mID {
		t.Error("incorrect message returned")
	}
}
