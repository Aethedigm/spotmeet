// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import "testing"

func TestMatch_Insert_DuplicateUser(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "test@test.com",
		Active:    1,
	}

	id, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	m := Match{
		User_A_ID:    id,
		User_B_ID:    id,
		PercentMatch: 100,
	}

	_, err = m.Insert(m)
	if err == nil {
		t.Error("failed to catch duplicate user")
	}
}

func TestMatch_MissingUserA(t *testing.T) {
	m := Match{
		User_B_ID: 1,
	}

	_, err := m.Insert(m)
	if err == nil {
		t.Error("failed to catch missing User_A")
	}
}

func TestMatch_MissingUserB(t *testing.T) {
	m := Match{
		User_A_ID: 1,
	}

	_, err := m.Insert(m)
	if err == nil {
		t.Error("failed to catch missing User_B")
	}
}

func TestMatch_Get_BadID(t *testing.T) {
	m := Match{}
	_, err := m.Get(1000)
	if err == nil {
		t.Error("failed to catch bad match id")
	}
}

func TestMatch_Get(t *testing.T) {

	u1 := &User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "test5@test.com",
		Active:    1,
	}

	u1ID, err := u1.Insert(*u1)
	if err != nil {
		t.Error(err)
	}

	u2 := &User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "test6@test.com",
		Active:    1,
	}

	u2ID, err := u2.Insert(*u2)
	if err != nil {
		t.Error(err)
	}

	m := Match{
		User_A_ID: u1ID,
		User_B_ID: u2ID,
	}

	mID, err := m.Insert(m)
	if err != nil {
		t.Error(err)
	}

	actual, err := m.Get(mID)
	if err != nil {
		t.Error(err)
	}

	if mID != actual.ID {
		t.Error("incorrect match returned")
	}
}
