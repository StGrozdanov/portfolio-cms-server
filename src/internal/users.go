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

// GetUserSkills gets the user skills info from the database (such as tech stack, soft skills and hobbies)
// and returns it
func GetUserSkills() (userSkills UserSkills, err error) {
	err = database.GetSingleRecord(
		&userSkills,
		`SELECT technology_stack ::JSONB -> 'techStack' AS tech_stack,
					   soft_skills ::JSONB -> 'softSkills'     AS soft_skills,
					   hobbies ::JSONB -> 'hobbies'            AS hobbies
				FROM users;`,
	)

	if err != nil {
		return UserSkills{}, err
	}
	return
}
