package data

import (
	up "github.com/upper/db/v4"
)

type Settings struct {
	ID                     int    `db:"id,omitempty"`
	UserID                 int    `db:"user_id" json:"user_id"`
	Distance               int    `db:"distance" json:"distance"`
	LookingFor             string `db:"looking_for" json:"looking_for"`
	MatchSensitivity       int    `db:"match_sensitivity" json:"match_sensitivity"`
	LikedArtistSensitivity int    `db:"liked_artist_sensitivity" json:"liked_artist_sensitivity"`
}

func (s *Settings) Table() string {
	return "settings"
}

func (s *Settings) GetAll() ([]*Settings, error) {
	collection := upper.Collection(s.Table())

	var all []*Settings

	res := collection.Find().OrderBy("id")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func (s *Settings) Get(id int) (*Settings, error) {
	var settings Settings
	collection := upper.Collection(s.Table())
	res := collection.Find(up.Cond{"id =": id})

	err := res.One(&settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func (s *Settings) GetByUserID(user_id int) (*Settings, error) {
	var settings Settings
	collection := upper.Collection(s.Table())
	res := collection.Find(up.Cond{"user_id =": user_id})

	err := res.One(&settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func (s *Settings) Delete(id int) error {
	collection := upper.Collection(s.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil

}

func (s *Settings) DeleteByUserID(user_id int) error {
	collection := upper.Collection(s.Table())
	res := collection.Find(up.Cond{"user_id =": user_id})
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil

}

func (s *Settings) Insert(settings Settings) (int, error) {

	// Default settings for all new users
	settings.MatchSensitivity = 50
	settings.LikedArtistSensitivity = 5
	settings.Distance = 50 // in miles
	settings.LookingFor = "friends"

	collection := upper.Collection(s.Table())
	res, err := collection.Insert(settings)
	if err != nil {
		return 0, err
	}

	id := getInsertID(res.ID())

	return id, nil
}

func (s *Settings) Update(settings Settings) error {
	collection := upper.Collection(s.Table())
	res := collection.Find(settings.ID)
	err := res.Update(settings)
	if err != nil {
		return err
	}
	return nil
}
