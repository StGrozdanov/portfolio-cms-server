package internal

import (
	"portfolio-cms-server/database"
)

// GetBasicUserInfo gets the basic user info from the database and returns it
func GetBasicUserInfo() (basicUserInfo BasicUserInfo, err error) {
	err = database.GetSingleRecord(
		&basicUserInfo,
		`SELECT email, cv_link, about_me, partners, carousel FROM users;`,
	)

	if err != nil {
		return BasicUserInfo{}, err
	}
	return
}
