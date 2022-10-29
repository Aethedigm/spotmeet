// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import "testing"

func TestLink_Table(t *testing.T) {
	s := Link{}
	table := s.Table()
	if table != "links" {
		t.Error("table name should be links")
	}
}

func TestLink_Insert(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "link_insert@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	u2 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "link_insert2@test.com",
		Active:    1,
	}

	uID2, err := u2.Insert(u2)
	if err != nil {
		t.Error(err)
	}

	a := Artist{
		SpotifyID: "test",
		Name:      "test",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}

	l := Link{
		User_A_ID: uID,
		User_B_ID: uID2,
		ArtistID:  aID,
	}

	lID, err := l.Insert(l)
	if err != nil {
		t.Error(err)
	}

	l.ID = lID

	if l.ID != lID {
		t.Error("link not inserted")
	}
}

func TestLink_Get_Bad_User_A_ID(t *testing.T) {
	ub := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "testub@test.com",
		Active:    1,
	}

	ubID, err := ub.Insert(ub)
	if err != nil {
		t.Error(err)
	}

	ub.ID = ubID

	a := Artist{
		SpotifyID: "testA",
		Name:      "testA",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}

	l := Link{
		User_A_ID: 0,
		User_B_ID: ubID,
		ArtistID:  aID,
	}

	_, err = l.Insert(l)
	if err == nil {
		t.Error("User A ID should not be 0")
	}
}

func TestLink_Get_Bad_User_B_ID(t *testing.T) {
	ua := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "testua@test.com",
		Active:    1,
	}

	uaID, err := ua.Insert(ua)
	if err != nil {
		t.Error(err)
	}

	ua.ID = uaID

	a := Artist{
		SpotifyID: "testA",
		Name:      "testA",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}

	l := Link{
		User_A_ID: uaID,
		User_B_ID: 0,
		ArtistID:  aID,
	}

	_, err = l.Insert(l)
	if err == nil {
		t.Error("User B ID should not be 0")
	}
}

func TestLink_Delete(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "link_delete@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	u2 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "link_delete2@test.com",
		Active:    1,
	}

	uID2, err := u2.Insert(u2)
	if err != nil {
		t.Error(err)
	}

	a := Artist{
		SpotifyID: "test_del",
		Name:      "test_del",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}

	l := Link{
		User_A_ID: uID,
		User_B_ID: uID2,
		ArtistID:  aID,
	}

	lID, err := l.Insert(l)
	if err != nil {
		t.Error(err)
	}

	l.ID = lID

	err = l.Delete(l.ID)
	if err != nil {
		t.Error(err)
	}

	_, err = l.Get(l.ID)
	if err == nil {
		t.Error("link not deleted")
	}
}

func TestLink_GetAll(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "l_getall@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	u2 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "l_getall2@test.com",
		Active:    1,
	}

	uID2, err := u2.Insert(u2)
	if err != nil {
		t.Error(err)
	}

	a := Artist{
		SpotifyID: "test_getall",
		Name:      "test_getall",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}

	l := Link{
		User_A_ID: uID,
		User_B_ID: uID2,
		ArtistID:  aID,
	}

	lID, err := l.Insert(l)
	if err != nil {
		t.Error(err)
	}

	l.ID = lID

	links, err := l.GetAll()
	if err != nil {
		t.Error(err)
	}

	if len(links) == 0 {
		t.Error("no links found")
	}
}

func TestLink_GetAllForOneUser(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "gafou_link@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	u2 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "gafou_link2@test.com",
		Active:    1,
	}

	uID2, err := u2.Insert(u2)
	if err != nil {
		t.Error(err)
	}

	a := Artist{
		SpotifyID: "test_gafou",
		Name:      "test_gafou",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}

	l := Link{
		User_A_ID: uID,
		User_B_ID: uID2,
		ArtistID:  aID,
	}

	lID, err := l.Insert(l)
	if err != nil {
		t.Error(err)
	}

	l.ID = lID

	links, err := l.GetAllForOneUser(uID)
	if err != nil {
		t.Error(err)
	}

	if len(links) == 0 {
		t.Error("no links found")
	}
}
