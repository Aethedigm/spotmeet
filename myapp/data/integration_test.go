// go:build integration

// run tests with this command: go test . --tags integration --count=1

package data

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "secret"
	dbName   = "celeritas_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var models Models
var testDB *sql.DB
var resource *dockertest.Resource
var pool *dockertest.Pool

func TestMain(m *testing.M) {
	os.Setenv("DATABASE_TYPE", "postgres")
	os.Setenv("UPPER_DB_LOG", "ERROR")

	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.4",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		// _ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to docker: %s", err)
	}

	err = createTables(testDB)
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	models = New(testDB)

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func compileMigrations(ext string) (string, error) {
	var res []string
	err := filepath.WalkDir("../migrations", func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		_, file := filepath.Split(s)

		if strings.HasSuffix(file, ext) {
			f, err := os.ReadFile(s)
			if err != nil {
				return err
			}

			res = append(res, string(f))
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return strings.Join(res, " "), nil
}

func createTables(db *sql.DB) error {
	stmt, err := compileMigrations(".up.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

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

func TestUser_GetAll(t *testing.T) {
	_, err := models.Users.GetAll()
	if err != nil {
		t.Error("failed to get user: ", err)
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

func TestUser_Delete(t *testing.T) {
	err := models.Users.Delete(1)
	if err != nil {
		t.Error("failed to delete user: ", err)
	}

	_, err = models.Users.Get(1)
	if err == nil {
		t.Error("retrieved user who was supposed to be deleted")
	}
}

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

func FindIDIn(id int, matches []*Match) bool {
	for _, m := range matches {
		if m.ID == id {
			return true
		}
	}
	return false
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

func TestArtist_GetAll(t *testing.T) {
	a := Artist{
		Name:      "test",
		SpotifyID: "123",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}
	a.ID = aID

	b := Artist{
		Name:      "test2",
		SpotifyID: "1234",
	}

	bID, err := b.Insert(b)
	if err != nil {
		t.Error(err)
	}
	b.ID = bID

	artists, err := a.GetAll()

	if ArtistInArray(a, artists) == false {
		t.Error("failed to return artist")
	}

	if ArtistInArray(b, artists) == false {
		t.Error("failed to return artist")
	}
}

func TestArtist_GetByName(t *testing.T) {
	a := Artist{
		Name:      "test3",
		SpotifyID: "123",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}
	a.ID = aID

	artist, err := a.GetByName(a.Name)
	if err != nil {
		t.Error(err)
	}

	if artist.ID != aID {
		t.Error("incorrect artist returned")
	}
}

func TestArtist_Get(t *testing.T) {
	a := Artist{
		Name:      "test4",
		SpotifyID: "123",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}
	a.ID = aID

	artist, err := a.Get(aID)
	if err != nil {
		t.Error(err)
	}

	if artist.ID != aID {
		t.Error("incorrect artist returned")
	}
}

func TestArtist_Update(t *testing.T) {
	a := Artist{
		Name:      "test5",
		SpotifyID: "123",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}
	a.ID = aID

	a.Name = "test6"

	err = a.Update(a)
	if err != nil {
		t.Error(err)
	}

	artist, err := a.Get(aID)
	if err != nil {
		t.Error(err)
	}

	if artist.Name != "test6" {
		t.Error("incorrect artist returned")
	}
}

func TestArtist_Delete(t *testing.T) {
	a := Artist{
		Name:      "test7",
		SpotifyID: "123",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}
	a.ID = aID

	err = a.Delete(aID)
	if err != nil {
		t.Error(err)
	}
}

func TestArtist_DeleteByName(t *testing.T) {
	a := Artist{
		Name:      "test8",
		SpotifyID: "123",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}
	a.ID = aID

	err = a.DeleteByName(a.Name)
	if err != nil {
		t.Error(err)
	}
}

func ArtistInArray(a Artist, arr []*Artist) bool {
	for _, v := range arr {
		if v.ID == a.ID {
			return true
		}
	}
	return false
}

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

func MessageInArray(m Message, arr []*Message) bool {
	for _, v := range arr {
		if v.ID == m.ID {
			return true
		}
	}
	return false
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

func TestProfile_Get(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "tpg1@test.com",
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

	profile, err := p.Get(pID)
	if err != nil {
		t.Error(err)
	}

	if profile.ID != pID {
		t.Error("incorrect profile returned")
	}
}

func TestProfile_GetByUserID(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "gbui1@test.com",
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

	profile, err := p.GetByUserID(uID)
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

func ProfileInArray(p Profile, arr []*Profile) bool {
	for _, v := range arr {
		if v.ID == p.ID {
			return true
		}
	}
	return false
}
