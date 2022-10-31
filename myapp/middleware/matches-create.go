package middleware

import (
	"log"
	"myapp/data"
)

// CreateMatches evaluates the given user's settings and musical tastes with all other users,
// and makes new match records in the Matches table for the given user with every other user that
// passes the evaluation. Returns the amount of new matches.
func CreateMatches(user int) (int, error) {

	// find matches for the user based only on location and looking-for settings
	potentialMatches := FilterByLocationAndLookingFor(user)

	// other match-filters (here soon) --------------------------------------

	// insert matches into the database
	var amountOfNewMatches = 0 // for counting how many new matches are created
	for _, match := range potentialMatches {
		m := data.Match{
			User_A_ID:    user,
			User_B_ID:    match,
			PercentMatch: 100,
			ArtistID:     0,
		}
		_, err := m.Insert(m)
		if err != nil {
			return 0, err
		}
		amountOfNewMatches += 1
	}
	return amountOfNewMatches, nil
}

// FilterByLocationAndLookingFor returns a slice of user_ids of potential-match users who
// fall within the distance and looking-for parameters given by the user and every other user.
func FilterByLocationAndLookingFor(userA int) []int {

	user := data.User{}

	//// get struct for given user_id
	//userAstruct, err := user.Get(userA)
	//if err != nil {
	//	log.Fatal(err)
	//	return []int{0}
	//}

	// get the distance setting of given user
	settings := data.Settings{}
	userAsettings, err := settings.GetByUserID(userA)
	if err != nil {
		log.Fatal(err)
	}
	maxDistanceForUserA := userAsettings.Distance

	// get structs for all users
	allUsers, err := user.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	// Check all users against given user, and save them in a slice if
	// users' settings are compatible with each other.
	var potentialMatches = []int{}
	for _, userB := range allUsers {
		if userA != userB.ID {
			// get userB settings
			userBsettings, err := settings.GetByUserID(userB.ID)
			if err != nil {
				log.Fatal(err)
			}
			maxDistanceForUserB := userAsettings.Distance

			// filter by distance
			distance, err := user.DistanceBetween(userA, userB.ID)
			if err != nil {
				log.Fatal(err)
			}
			if distance <= maxDistanceForUserA && distance <= maxDistanceForUserB {
				// filter by looking-for
				if userAsettings.LookingFor == userBsettings.LookingFor {
					// append userB.ID to the slice
					potentialMatches = append(potentialMatches, userB.ID)
				}
			}
		}
	}

	return potentialMatches
}
