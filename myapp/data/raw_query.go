package data

import (
	"fmt"
	"strconv"
	"time"
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

func (r *RawQuery) ThreadPreviewQuery(userID int, otherUserID int) (string, string, string, error) {

	// query to find content of latest message, and the time for that message
	q1 := `select m.content, m.created_at
	from messages m 
	where (m.user_id = ` + strconv.Itoa(userID) + ` and m.match_id = ` + strconv.Itoa(otherUserID) + `)
	or (m.match_id = ` + strconv.Itoa(userID) + ` and m.user_id = ` + strconv.Itoa(otherUserID) + `)
	order by m.created_at::DATE DESC, m.created_at DESC;`

	q1rows, err := upper.SQL().Query(q1)
	if err != nil {
		fmt.Println("problem with query within func ThreadPreviewQuery", q1rows, err)
		return "", "", "", err
	}

	// query to get other-user's profile image url
	q2 := `select p.profile_image_url
	from profiles p 
	where p.user_id = ` + strconv.Itoa(otherUserID) + `;`

	q2rows, err := upper.SQL().Query(q2)
	if err != nil {
		fmt.Println("problem with query within func ThreadPreviewQuery", q2rows, err)
		return "", "", "", err
	}

	// create containers to return
	var LatestMessagePreview string
	var LatestMessageTimeSent time.Time
	var OtherUsersImage string

	// pull first query's top record into topRecord struct
	q1rows.Next()
	q1rows.Scan(&LatestMessagePreview, &LatestMessageTimeSent)

	// get the final variables to return
	if len(LatestMessagePreview) > 35 {
		LatestMessagePreview = LatestMessagePreview[:35] + " . . ."
	}
	strLatestMessageTimeSent := LatestMessageTimeSent.Format(time.Kitchen)

	q2rows.Next()
	q2rows.Scan(&OtherUsersImage)

	return LatestMessagePreview, strLatestMessageTimeSent, OtherUsersImage, nil
}
