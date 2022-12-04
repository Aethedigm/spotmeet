package data

import (
	"fmt"
	"strconv"

	up "github.com/upper/db/v4"
)

type Settings struct {
	ID                     int     `db:"id,omitempty"`
	UserID                 int     `db:"user_id" json:"user_id"`
	Distance               int     `db:"distance" json:"distance"`
	LookingFor             string  `db:"looking_for" json:"looking_for"`
	MatchSensitivity       int     `db:"match_sensitivity" json:"match_sensitivity"`
	LikedArtistSensitivity int     `db:"liked_artist_sensitivity" json:"liked_artist_sensitivity"`
	LatMin                 float64 `db:"lat_min" json:"lat_min"`
	LatMax                 float64 `db:"lat_max" json:"lat_max"`
	LongMin                float64 `db:"long_min" json:"long_min"`
	LongMax                float64 `db:"long_max" json:"long_max"`
}

func (s *Settings) MatchSensitivityString() string {
	switch s.MatchSensitivity {
	case 1:
		return "Low"
	case 2:
		return "Medium"
	case 3:
		return "High"
	}

	return "Low"
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
	currentUserSettings, err := s.GetByUserID(settings.UserID)
	if err != nil {
		return err
	}

	// wipe unlinked matches if looking-for and/or distance is changed
	if settings.Distance != currentUserSettings.Distance ||
		settings.LookingFor != currentUserSettings.LookingFor {
		userID := settings.UserID
		q := `delete 
			from matches 
			where user_a_id = ` + strconv.Itoa(userID) + ` and complete = false ` +
			`or user_b_id = ` + strconv.Itoa(userID) + ` and complete = false;`

		rows, err := upper.SQL().Query(q)
		if err != nil {
			fmt.Println("problem with query", rows, err)
			return err
		}
	}

	collection := upper.Collection(s.Table())
	res := collection.Find(settings.ID)
	err = res.Update(settings)
	if err != nil {
		return err
	}
	return nil
}
