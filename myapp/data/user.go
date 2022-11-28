package data

import (
	"errors"
	"time"

	up "github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `db:"id,omitempty"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Email     string    `db:"email" json:"email"`
	Active    int       `db:"user_active"`
	Password  string    `db:"password" json:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Token     Token     `db:"-"`
	Latitude  float64   `db:"lat" json:"lat"`
	Longitude float64   `db:"long" json:"long"`
}

func (u *User) Table() string {
	return "users"
}

func (u *User) GetAll() ([]*User, error) {
	collection := upper.Collection(u.Table())

	var all []*User

	res := collection.Find().OrderBy("last_name")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	var theUser User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"email": email})
	err := res.One(&theUser)
	if err != nil {
		return nil, err
	}

	var token Token
	collection = upper.Collection(token.Table())
	res = collection.Find(up.Cond{"user_id": theUser.ID, "expiry >": time.Now()}).OrderBy("created_at desc")
	err = res.One(&token)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}

	theUser.Token = token

	return &theUser, nil
}

func (u *User) Get(id int) (*User, error) {
	var theUser User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"id =": id})

	err := res.One(&theUser)
	if err != nil {
		return nil, err
	}

	var token Token
	collection = upper.Collection(token.Table())
	res = collection.Find(up.Cond{"user_id =": theUser.ID, "expiry >": time.Now()}).OrderBy("created_at desc")
	err = res.One(&token)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}

	theUser.Token = token

	return &theUser, nil
}

func (u *User) Update(theUser User) error {
	theUser.UpdatedAt = time.Now()
	collection := upper.Collection(u.Table())
	res := collection.Find(theUser.ID)
	err := res.Update(&theUser)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete(id int) error {
	collection := upper.Collection(u.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil

}

func (u *User) Insert(theUser User) (int, error) {
	collection := upper.Collection(u.Table())

	newHash, err := bcrypt.GenerateFromPassword([]byte(theUser.Password), 12)
	if err != nil {
		return 0, err
	}

	theUser.Active = 1
	theUser.CreatedAt = time.Now()
	theUser.UpdatedAt = time.Now()
	theUser.Password = string(newHash)

	res, err := collection.Insert(theUser)
	if err != nil {
		return 0, err
	}

	uID := getInsertID(res.ID())

	return uID, nil
}

func (u *User) ResetPassword(id int, password string) error {
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	theUser, err := u.Get(id)
	if err != nil {
		return err
	}

	theUser.Password = string(newHash)

	err = u.Update(*theUser)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			// some kind of error occurred
			return false, err
		}
	}

	return true, nil
}
