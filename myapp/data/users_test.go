// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import "testing"

func TestUser_Table(t *testing.T) {
	s := models.Users.Table()
	if s != "users" {
		t.Error("wrong table name returned: ", s)
	}
}

func TestUser_Insert(t *testing.T) {
	u := User{
		FirstName: "test",
		LastName:  "test",
		Email:     "user_insert@test.com",
	}

	id, err := u.Insert(u)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	if id == 0 {
		t.Error("0 returned as id after insert")
	}
}

func TestUser_Insert_Duplicate(t *testing.T) {
	u1 := User{
		FirstName: "test",
		LastName:  "test",
		Email:     "user_insert_duplicate",
	}

	id, err := u1.Insert(u1)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	if id == 0 {
		t.Error("0 returned as id after insert")
	}

	_, err = u1.Insert(u1)
	if err == nil {
		t.Error("no error returned when inserting duplicate")
	}
}

func TestUser_Get(t *testing.T) {
	u, err := models.Users.Get(1)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	if u.ID == 0 {
		t.Error("id of returned user is 0: ", err)
	}
}

func TestUser_GetByEmail(t *testing.T) {

	u1 := User{
		FirstName: "test",
		LastName:  "test",
		Email:     "user_getbyemail@test.com",
	}

	_, err := u1.Insert(u1)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	u, err := models.Users.GetByEmail(u1.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	if u.ID == 0 {
		t.Error("id of returned user is 0: ", err)
	}
}

func TestUser_Update(t *testing.T) {
	u, err := models.Users.Get(1)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	u.LastName = "Smith"
	err = u.Update(*u)
	if err != nil {
		t.Error("failed to update user: ", err)
	}

	u, err = models.Users.Get(1)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	if u.LastName != "Smith" {
		t.Error("last name not updated in database")
	}
}

func TestUser_PasswordMatches(t *testing.T) {

	u1 := User{
		FirstName: "test",
		LastName:  "test",
		Email:     "user_passwordmatches@test.com",
		Password:  "password",
	}

	u1ID, err := u1.Insert(u1)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	u, err := models.Users.Get(u1ID)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	matches, err := u.PasswordMatches("password")
	if err != nil {
		t.Error("error checking match: ", err)
	}

	if !matches {
		t.Error("password doesn't match when it should")
	}

	matches, err = u.PasswordMatches("123")
	if err != nil {
		t.Error("error checking match: ", err)
	}

	if matches {
		t.Error("password matches when it should not")
	}
}

func TestUser_ResetPassword(t *testing.T) {

	u := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "resetpassword@test.com",
	}

	id, err := models.Users.Insert(u)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	err = models.Users.ResetPassword(id, "new_password")
	if err != nil {
		t.Error("error resetting password: ", err)
	}

	err = models.Users.ResetPassword(1000, "new_password")
	if err == nil {
		t.Error("did not get an error when trying to reset password for non-existent user")
	}
}
