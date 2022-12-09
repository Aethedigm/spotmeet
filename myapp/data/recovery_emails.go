// Spotmeet - (Capstone Team E)
// 2022 Stephen Sumpter, John Neumeier,
// Zach Kods, Landon Wilson
package data

import up "github.com/upper/db/v4"

type RecoveryEmail struct {
	ID     int `json:"id" db:"id,omitempty"`
	UserID int `json:"userid" db:"userid"`
}

func (re *RecoveryEmail) Table() string {
	return "recovery_emails"
}

func (re *RecoveryEmail) Get(id int) (*RecoveryEmail, error) {
	collection := upper.Collection(re.Table())

	var rec RecoveryEmail

	res := collection.Find(up.Cond{"id": id})
	err := res.One(&rec)
	if err != nil {
		return nil, err
	}

	return &rec, nil
}

func (re *RecoveryEmail) Insert(rec RecoveryEmail) (int, error) {
	collection := upper.Collection(re.Table())
	res, err := collection.Insert(rec)
	if err != nil {
		return 0, err
	}
	id := getInsertID(res.ID())

	return id, nil
}

func (re *RecoveryEmail) DeleteAllForUser(userID int) error {
	collection := upper.Collection(re.Table())
	res := collection.Find(up.Cond{"userid": userID})
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}
