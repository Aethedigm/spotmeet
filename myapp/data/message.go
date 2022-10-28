package data

import (
	"errors"
	"time"

	up "github.com/upper/db/v4"
)

// Message is the type for a message
type Message struct {
	ID        int       `db:"id,omitempty"`
	UserID    int       `db:"user_id" json:"user_id"`
	MatchID   int       `db:"match_id" json:"match_id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at"`
}

// Table returns the table name associated with this model in the database
func (m *Message) Table() string {
	return "messages"
}

// GetAll returns a slice of all messages.
func (m *Message) GetAll() ([]*Message, error) {
	collection := upper.Collection(m.Table())

	var all []*Message

	res := collection.Find().OrderBy("id")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// GetAllForOneMatch returns a slice of all messages for a match.
func (m *Message) GetAllForOneMatch(matchID int) ([]*Message, error) {

	var all []*Message
	var tmp []*Message

	collection := upper.Collection(m.Table())
	res1 := collection.Find(up.Cond{"match_id": matchID})
	res2 := collection.Find(up.Cond{"user_id": matchID})

	err := res1.All(&all)
	if err != nil {
		return nil, err
	}

	err = res2.All(&tmp)
	if err != nil {
		return nil, err
	}

	all = append(all, tmp...)

	return all, nil
}

// Get gets one message by id
func (m *Message) Get(id int) (*Message, error) {
	var themessage Message
	collection := upper.Collection(m.Table())
	res := collection.Find(up.Cond{"id": id})

	err := res.One(&themessage)
	if err != nil {
		return nil, err
	}

	return &themessage, nil
}

// Delete deletes a message by id
func (m *Message) Delete(id int) error {
	collection := upper.Collection(m.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a new message, and returns the newly inserted message's id
func (m *Message) Insert(themessage Message) (int, error) {

	if themessage.UserID == 0 || themessage.MatchID == 0 {
		return 0, errors.New("UserID and MatchID must be set")
	}

	themessage.CreatedAt = time.Now()

	collection := upper.Collection(m.Table())
	res, err := collection.Insert(themessage)
	if err != nil {
		return 0, err
	}

	id := getInsertID(res.ID())

	return id, nil
}
