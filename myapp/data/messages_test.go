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

func TestMessage_GetAllForMatch(t *testing.T) {
	u1 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "gafm1@test.com",
		Active:    1,
	}

	u1ID, err := u1.Insert(u1)
	if err != nil {
		t.Error(err)
	}

	u2 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "gafm2@test.com",
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

	m2 := Message{
		UserID:  u2ID,
		MatchID: u1ID,
		Content: "test2",
	}

	m2ID, err := m2.Insert(m2)
	if err != nil {
		t.Error(err)
	}

	m2.ID = m2ID

	messages, err := m.GetAllForOneMatch(u1ID)
	if err != nil {
		t.Error(err)
	}

	if len(messages) != 2 {
		t.Error("incorrect number of messages returned")
	}

	if !MessageInArray(m, messages) {
		t.Error("incorrect message returned")
	}
}

func TestMessage_GetAll(t *testing.T) {
	u1 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "ga1@test.com",
		Active:    1,
	}

	u1ID, err := u1.Insert(u1)
	if err != nil {
		t.Error(err)
	}

	u2 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "ga2@test.com",
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

	m2 := Message{
		UserID:  u2ID,
		MatchID: u1ID,
		Content: "test2",
	}

	m2ID, err := m2.Insert(m2)
	if err != nil {
		t.Error(err)
	}

	m2.ID = m2ID

	messages, err := m.GetAll()
	if err != nil {
		t.Error(err)
	}

	if !MessageInArray(m, messages) {
		t.Error("incorrect message returned")
	}
}

func TestMessage_Delete(t *testing.T) {
	u1 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "d1@test.com",
		Active:    1,
	}

	u1ID, err := u1.Insert(u1)
	if err != nil {
		t.Error(err)
	}

	u2 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "d2@test.com",
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

	err = m.Delete(mID)
	if err != nil {
		t.Error(err)
	}

	_, err = m.Get(mID)
	if err == nil {
		t.Error("message not deleted")
	}
}
