// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import (
	"net/http"
	"testing"
	"time"
)

func TestToken_Table(t *testing.T) {
	s := models.Tokens.Table()
	if s != "tokens" {
		t.Error("wrong table name returned for tokens")
	}
}

func TestToken_GenerateToken(t *testing.T) {

	u := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "generatetoken@test.com",
	}

	id, err := models.Users.Insert(u)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	_, err = models.Tokens.GenerateToken(id, time.Hour*24*365)
	if err != nil {
		t.Error("error generating token: ", err)
	}
}

func TestToken_Insert(t *testing.T) {
	u := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "tokeninsert@test.com",
	}

	id, err := models.Users.Insert(u)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	token, err := models.Tokens.GenerateToken(id, time.Hour*24*365)
	if err != nil {
		t.Error("error generating token: ", err)
	}

	err = models.Tokens.Insert(*token, u)
	if err != nil {
		t.Error("error insering token: ", err)
	}
}
func TestToken_GetUserForToken(t *testing.T) {
	token := "abc"
	_, err := models.Tokens.GetUserForToken(token)
	if err == nil {
		t.Error("error expected but not recieved when getting user with a bad token")
	}

	u1 := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "getuserfortoken@test.com",
	}

	u1ID, err := u1.Insert(u1)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	tok, err := models.Tokens.GenerateToken(u1ID, time.Hour*24*365)
	if err != nil {
		t.Error("error generating token: ", err)
	}

	err = tok.Insert(*tok, u1)
	if err != nil {
		t.Error("error inserting token: ", err)
	}

	u, err := u1.GetByEmail(u1.Email)
	if err != nil {
		t.Error("failed to get user")
	}

	_, err = models.Tokens.GetUserForToken(u.Token.PlainText)
	if err != nil {
		t.Error("failed to get user with valid token: ", err)
	}
}

func TestToken_GetTokensForUser(t *testing.T) {
	tokens, err := models.Tokens.GetTokensForUser(1)
	if err != nil {
		t.Error(err)
	}

	if len(tokens) > 0 {
		t.Error("tokens returned for non-existent user")
	}
}

func TestToken_GetByToken(t *testing.T) {

	u1 := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "token_getbytoken@test.com",
	}

	_, err := u1.Insert(u1)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	u, err := models.Users.GetByEmail(u1.Email)
	if err != nil {
		t.Error("failed to get user")
	}

	tok, err := models.Tokens.GenerateToken(u.ID, time.Hour*24*365)
	if err != nil {
		t.Error("error generating token: ", err)
	}

	err = tok.Insert(*tok, u1)
	if err != nil {
		t.Error("error inserting token: ", err)
	}

	_, err = models.Tokens.GetByToken(tok.PlainText)
	if err != nil {
		t.Error("error getting token by token: ", err)
	}

	_, err = models.Tokens.GetByToken("123")
	if err == nil {
		t.Error("no error getting non-existing token by token: ", err)
	}
}

func TestToken_AuthenticateToken(t *testing.T) {

	dummyUser := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "token_authenticatetoken@test.com",
	}

	dID, err := dummyUser.Insert(dummyUser)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	dummyUser.ID = dID

	tok, err := models.Tokens.GenerateAndInsert(dummyUser, time.Hour*24*365)
	if err != nil {
		t.Error("error generating token: ", err)
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tok.PlainText)

	_, err = models.Tokens.AuthenticateToken(req)
	if err != nil {
		t.Error("error authenticating token: ", err)
	}
}

func TestToken_Delete(t *testing.T) {

	dummyUser := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "token_delete@test.com",
	}

	dID, err := dummyUser.Insert(dummyUser)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	dummyUser.ID = dID

	tok, err := models.Tokens.GenerateAndInsert(dummyUser, time.Hour*24*365)
	if err != nil {
		t.Error("error generating token: ", err)
	}

	err = models.Tokens.DeleteByToken(tok.PlainText)
	if err != nil {
		t.Error("error deleting token: ", err)
	}
}

func TestToken_ExpiredToken(t *testing.T) {
	dummyUser := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "token_expiredtoken@test.com",
	}

	dID, err := dummyUser.Insert(dummyUser)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	dummyUser.ID = dID

	u, err := models.Users.GetByEmail(dummyUser.Email)
	if err != nil {
		t.Error(err)
	}

	token, err := models.Tokens.GenerateToken(u.ID, -time.Hour)
	if err != nil {
		t.Error(err)
	}

	err = models.Tokens.Insert(*token, *u)
	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer "+token.PlainText)

	_, err = models.Tokens.AuthenticateToken(req)
	if err == nil {
		t.Error("failed to catch expired token")
	}

}

func TestToken_BadHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	_, err := models.Tokens.AuthenticateToken(req)
	if err == nil {
		t.Error("failed to catch missing auth header")
	}

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "abc")
	_, err = models.Tokens.AuthenticateToken(req)
	if err == nil {
		t.Error("failed to catch bad auth header")
	}

	newUser := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "you@there.com",
		Active:    1,
		Password:  "abc",
	}

	id, err := models.Users.Insert(newUser)
	if err != nil {
		t.Error(err)
	}

	token, err := models.Tokens.GenerateToken(id, 1*time.Hour)
	if err != nil {
		t.Error(err)
	}

	err = models.Tokens.Insert(*token, newUser)
	if err != nil {
		t.Error(err)
	}

	err = models.Users.Delete(id)
	if err != nil {
		t.Error(err)
	}

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer "+token.PlainText)
	_, err = models.Tokens.AuthenticateToken(req)
	if err == nil {
		t.Error("failed to catch token for deleted user")
	}

}

func TestToken_DeleteNonExistentToken(t *testing.T) {
	err := models.Tokens.DeleteByToken("abc")
	if err != nil {
		t.Error("error deleting token")
	}
}

func TestToken_ValidToken(t *testing.T) {
	dummyUser := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "token_validtoken@test.com",
	}

	dID, err := dummyUser.Insert(dummyUser)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	dummyUser.ID = dID

	u, err := models.Users.GetByEmail(dummyUser.Email)
	if err != nil {
		t.Error(err)
	}

	newToken, err := models.Tokens.GenerateToken(u.ID, 24*time.Hour)
	if err != nil {
		t.Error(err)
	}

	err = models.Tokens.Insert(*newToken, *u)
	if err != nil {
		t.Error(err)
	}

	okay, _ := models.Tokens.ValidToken(newToken.PlainText)
	if !okay {
		t.Error("valid token reported as invalid")
	}

	okay, _ = models.Tokens.ValidToken("abc")
	if okay {
		t.Error("invalid token reported as valid")
	}

	u, err = models.Users.GetByEmail(dummyUser.Email)
	if err != nil {
		t.Error(err)
	}

	err = models.Tokens.Delete(u.Token.ID)
	if err != nil {
		t.Error(err)
	}

	okay, err = models.Tokens.ValidToken(u.Token.PlainText)
	if err == nil {
		t.Error(err)
	}
	if okay {
		t.Error("no error reported when validating non-existent token")
	}
}
