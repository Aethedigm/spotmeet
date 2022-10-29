// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import "testing"

func TestProfile_GetByUserID(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "templast",
		Email:     "gbuid@test.com",
	}

	uID, err := models.Users.Insert(u)
	if err != nil {
		t.Error(err)
	}

	p := Profile{
		UserID: uID,
	}

	pID, err := p.Insert(p)
	if err != nil {
		t.Error(err)
	}

	p.ID = pID

	profile, err := p.GetByUserID(uID)
	if err != nil {
		t.Error(err)
	}

	if profile.ID != pID {
		t.Error("incorrect profile returned")
	}
}

func TestProfile_Get(t *testing.T) {
	u := &User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "tpg123@test.com",
		Active:    1,
	}

	uID, err := u.Insert(*u)
	if err != nil {
		t.Error(err)
	}

	p := Profile{
		UserID: uID,
	}

	pID, err := p.Insert(p)
	if err != nil {
		t.Error(err)
	}

	p.ID = pID

	profile, err := p.Get(pID)
	if err != nil {
		t.Error(err)
	}

	if profile.ID != pID {
		t.Error("incorrect profile returned")
	}
}

func TestProfile_Update(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	p := Profile{
		UserID: uID,
	}

	pID, err := p.Insert(p)
	if err != nil {
		t.Error(err)
	}

	p.ID = pID

	p.Description = "test"

	err = p.Update(p)
	if err != nil {
		t.Error(err)
	}

	profile, err := p.Get(pID)
	if err != nil {
		t.Error(err)
	}

	if profile.Description != "test" {
		t.Error("incorrect profile returned")
	}
}

func TestProfile_Delete(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "del@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	p := Profile{
		UserID: uID,
	}

	pID, err := p.Insert(p)
	if err != nil {
		t.Error(err)
	}

	p.ID = pID

	err = p.Delete(pID)
	if err != nil {
		t.Error(err)
	}

	_, err = p.Get(pID)
	if err == nil {
		t.Error("profile not deleted")
	}
}

func TestProfile_GetAll(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "gall@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	p := Profile{
		UserID: uID,
	}

	pID, err := p.Insert(p)
	if err != nil {
		t.Error(err)
	}

	p.ID = pID

	profiles, err := p.GetAll()
	if err != nil {
		t.Error(err)
	}

	if !ProfileInArray(p, profiles) {
		t.Error("missing profile")
	}
}
