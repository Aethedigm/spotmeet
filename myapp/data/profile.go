package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// Profile is a struct that holds profile data for a user
type Profile struct {
	ID          int       `db:"id,omitempty"`
	UserID      int       `db:"user_id" json:"user_id"`
	Description string    `db:"description" json:"description"`
	ImageURL    string    `db:"profile_image_url" json:"profile_image_url"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Table returns the table name associated with this model in the database
func (p *Profile) Table() string {
	return "profiles"
}

// GetAll returns a slice of all profiles
func (p *Profile) GetAll() ([]Profile, error) {
	collection := upper.Collection(p.Table())

	var all []Profile

	res := collection.Find().OrderBy("user_id")
	err := res.All(all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// Get gets one user by profile id
func (p *Profile) Get(id int) (*Profile, error) {
	var profile Profile
	collection := upper.Collection(p.Table())
	res := collection.Find(up.Cond{"id =": id})

	err := res.One(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

// Get gets one user by user id
func (p *Profile) GetByUserID(user_id int) (*Profile, error) {
	var profile Profile
	collection := upper.Collection(p.Table())
	res := collection.Find(up.Cond{"user_id =": user_id})

	err := res.One(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

// Update updates a profile record in the database
func (p *Profile) Update(profile Profile) error {
	profile.UpdatedAt = time.Now()
	collection := upper.Collection(p.Table())
	res := collection.Find(profile.ID)
	err := res.Update(&profile)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a profile by id
func (p *Profile) Delete(id int) error {
	collection := upper.Collection(p.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a profile by user id
func (p *Profile) DeleteByUserID(user_id int) error {
	collection := upper.Collection(p.Table())
	res := collection.Find(up.Cond{"user_id =": user_id})
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a new profile, and returns the newly inserted profile id
func (p *Profile) Insert(profile Profile) (int, error) {

	// enable this once default image is found.
	// ensure a default image is set if one was not set while user was creating profile
	//if profile.ImageURL == nil {
	//	profile.ImageURL = //////// some default image in public folder
	//}

	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	collection := upper.Collection(p.Table())
	res, err := collection.Insert(profile)
	if err != nil {
		return 0, err
	}

	id := getInsertID(res.ID())

	return id, nil
}
