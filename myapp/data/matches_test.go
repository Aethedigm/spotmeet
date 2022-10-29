// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import "testing"

func TestMatch_GetAll(t *testing.T) {
	u1 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "testmatch_getall1@test.com",
		Active:    1,
	}

	u2 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "testmatch_getall2@test.com",
		Active:    1,
	}

	u1ID, err := u1.Insert(u1)
	if err != nil {
		t.Error(err)
	}

	u2ID, err := u2.Insert(u2)
	if err != nil {
		t.Error(err)
	}

	m1 := Match{
		User_A_ID: u1ID,
		User_B_ID: u2ID,
	}

	m2 := Match{
		User_A_ID: u2ID,
		User_B_ID: u1ID,
	}

	m1ID, err := m1.Insert(m1)
	if err != nil {
		t.Error(err)
	}

	m2ID, err := m2.Insert(m2)
	if err != nil {
		t.Error(err)
	}

	actual, err := m1.GetAll()
	if err != nil {
		t.Error(err)
	}

	expectedIDs := []int{m1ID, m2ID}
	for _, id := range expectedIDs {
		if !FindIDIn(id, actual) {
			t.Error("expected match not found")
		}
	}
}

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

func TestMatch_GetAllForOneUser(t *testing.T) {
	ua := &User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "test2@test.com",
		Active:    1,
	}

	ub := &User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "test3@test.com",
		Active:    1,
	}

	uaID, err := ua.Insert(*ua)
	if err != nil {
		t.Error(err)
	}

	ubID, err := ub.Insert(*ub)
	if err != nil {
		t.Error(err)
	}

	m := Match{
		User_A_ID:    uaID,
		User_B_ID:    ubID,
		PercentMatch: 100,
	}

	mID, err := m.Insert(m)
	if err != nil {
		t.Error(err)
	}

	matches, err := m.GetAllForOneUser(uaID)
	if err != nil {
		t.Error(err)
	}

	if len(matches) < 1 {
		t.Error("no matches returned")
	}

	if mID != matches[0].ID {
		t.Error("incorrect match returned")
	}
}

func TestMatch_GetAllForOneUser_NoMatches(t *testing.T) {
	ua := &User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "test4@test.com",
		Active:    1,
	}

	uaID, err := ua.Insert(*ua)
	if err != nil {
		t.Error(err)
	}

	m := Match{
		User_A_ID:    uaID,
		User_B_ID:    uaID,
		PercentMatch: 100,
	}

	matches, err := m.GetAllForOneUser(uaID)
	if err != nil {
		t.Error(err)
	}

	if len(matches) > 0 {
		t.Error("matches returned")
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

func TestMatch_Delete(t *testing.T) {
	u1 := &User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "test7@test.com",
		Active:    1,
	}

	u1ID, err := u1.Insert(*u1)
	if err != nil {
		t.Error(err)
	}

	u2 := &User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "test8@test.com",
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

	err = m.Delete(mID)
	if err != nil {
		t.Error("Failed to delete match ", err)
	}
}
