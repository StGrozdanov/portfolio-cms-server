package users

import (
	"portfolio-cms-server/database"
)

// GetBasicInfo gets the basic user info from the database and returns it
func GetBasicInfo() (basicUserInfo BasicUserInfo, err error) {
	err = database.GetSingleRecord(
		&basicUserInfo,
		`SELECT email, cv_link, about_me, partners, carousel FROM users;`,
	)

	if err != nil {
		return BasicUserInfo{}, err
	}
	return
}

// GetSkills gets the user skills info from the database (such as tech stack, soft skills and hobbies)
// and returns it
func GetSkills() (userSkills UserSkills, err error) {
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

// GetJobsAndProjects gets the user jobs and projects info from the database
func GetJobsAndProjects() (jobsAndProjects JobsAndProjects, err error) {
	err = database.GetSingleRecord(&jobsAndProjects, `SELECT jobs, projects FROM users;`)
	if err != nil {
		return JobsAndProjects{}, err
	}
	return
}

// GetSocials gets the user jobs and projects info from the database
func GetSocials() (socials Socials, err error) {
	err = database.GetSingleRecord(&socials, `SELECT social_media FROM users;`)
	if err != nil {
		return Socials{}, err
	}
	return
}