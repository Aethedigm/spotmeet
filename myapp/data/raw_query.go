package data

import (
	"fmt"
)

type RawQuery struct {
	ID int
}

func (r *RawQuery) MatchQuery(user User, settings Settings) ([]int, error) {
	q := `select u.id
	from users u 
	where u.id in (
		select u.id
		from users u
		join settings s 
		on u.id = s.user_id 
		where 
			u.lat >= ` + fmt.Sprintf("%f", settings.LatMin) + `
		and u.lat <= ` + fmt.Sprintf("%f", settings.LatMax) + `
		and u.long <= ` + fmt.Sprintf("%f", settings.LongMax) + `
		and u.long >= ` + fmt.Sprintf("%f", settings.LongMin) + `
		and s.lat_max >= ` + fmt.Sprintf("%f", user.Latitude) + `
		and s.lat_min <= ` + fmt.Sprintf("%f", user.Latitude) + `
		and s.long_max >= ` + fmt.Sprintf("%f", user.Longitude) + `
		and s.long_min <= ` + fmt.Sprintf("%f", user.Longitude) + `
		and u.id != ` + fmt.Sprintf("%d", user.ID) + `
		and s.looking_for = '` + settings.LookingFor + `'
		and u.id not in (
			select m.user_a_id
			from matches m
			where m.user_b_id = ` + fmt.Sprintf("%d", user.ID) + `
		)
		and u.id not in (
			select m.user_b_id
			from matches m
			where m.user_a_id = ` + fmt.Sprintf("%d", user.ID) + `
		)
	);`

	rows, err := upper.SQL().Query(q)
	if err != nil {
		fmt.Println("problem with query", rows, err)
		return nil, err
	}

	var userIDs []int
	for rows.Next() {
		var u int
		err := rows.Scan(&u)
		if err != nil {
			fmt.Println("problem with filling users", u, err)
			return nil, err
		}
		userIDs = append(userIDs, u)
	}

	return userIDs, nil
}
