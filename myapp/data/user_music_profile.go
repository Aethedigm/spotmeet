package data

import (
	"time"

	up "github.com/upper/db/v4"
)

type UserMusicProfile struct {
	ID               int       `db:"id,omitempty" json:"id"`
	UserID           int       `db:"user_id" json:"user_id"`
	Loudness         float64   `db:"loudness" json:"loudness"`
	Tempo            float64   `db:"tempo" json:"tempo"`
	TimeSig          int       `db:"time_sig" json:"time_sig"`
	Acousticness     float32   `db:"acousticness" json:"acousticness"`
	Danceability     float32   `db:"danceability" json:"danceability"`
	Energy           float32   `db:"energy" json:"energy"`
	Instrumentalness float32   `db:"instrumentalness" json:"instrumentalness"`
	Mode             int       `db:"mode" json:"mode"`
	Speechiness      float32   `db:"speechiness" json:"speechiness"`
	Valence          float32   `db:"valence" json:"valence"`
	UpdatedAt        time.Time `db:"updated_at" json:"update_at"`
}

func (u *UserMusicProfile) Table() string {
	return "user_music_profile"
}

func (u *UserMusicProfile) Get(id int) (*UserMusicProfile, error) {
	var ump UserMusicProfile

	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"id": id})
	err := res.One(&ump)
	if err != nil {
		return nil, err
	}

	return &ump, nil
}

func (u *UserMusicProfile) GetByUser(user User) (*UserMusicProfile, error) {
	var ump UserMusicProfile

	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"user_id": user.ID})
	err := res.One(&ump)
	if err != nil {
		return nil, err
	}

	return &ump, nil
}

func (u *UserMusicProfile) GetByUserID(userID int) (*UserMusicProfile, error) {
	var ump UserMusicProfile

	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"user_id": userID})

	exists, err := res.Exists()
	if err != nil {
		return nil, err
	}

	if exists {
		err := res.One(&ump)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &ump, nil
}

func (u *UserMusicProfile) Insert(ump *UserMusicProfile) (int, error) {
	collection := upper.Collection(u.Table())

	ump.UpdatedAt = time.Now()

	res, err := collection.Insert(ump)
	if err != nil {
		return 0, err
	}

	id := getInsertID(res.ID())

	return id, nil
}

func (u *UserMusicProfile) Update(ump UserMusicProfile) (int, error) {
	ump.UpdatedAt = time.Now()
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{"user_id": ump.UserID})
	err := res.Update(ump)
	if err != nil {
		return 0, err
	}

	return ump.ID, nil
}
