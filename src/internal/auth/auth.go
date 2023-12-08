package auth

import (
	"golang.org/x/crypto/bcrypt"
	"portfolio-cms-server/database"
	"portfolio-cms-server/utils"
)

// Login accepts username and password validates it and if such user exists - returns JWT auth token
func Login(loginData UserAuthData) (authToken string, err error) {
	userData := UserAuthData{}

	err = database.GetSingleRecordNamedQuery(
		&userData,
		`SELECT nickname, password FROM users WHERE nickname = :username`,
		loginData,
	)
	if err != nil {
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(loginData.Password)); err != nil {
		return
	}

	return utils.GenerateJWT()
}
