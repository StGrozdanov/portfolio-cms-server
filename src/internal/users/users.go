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

// UpdateInfo updates the basic user info in the database and returns it
func UpdateInfo(request BasicUserInfo) (basicUserInfo BasicUserInfo, err error) {
	err = database.GetSingleRecordNamedQuery(
		&basicUserInfo,
		`UPDATE users
				SET about_me = :about_me,
					cv_link  = :cv_link,
					email    = :email,
					partners = :partners,
					carousel = :carousel
				WHERE users.id = 1
				RETURNING*;`,
		request,
	)

	if err != nil {
		return
	}
	return
}

// UpdateSkills updates the user skills in the database and returns it
func UpdateSkills(request UserSkills) (skills UserSkills, err error) {
	err = database.GetSingleRecordNamedQuery(
		&skills,
		`UPDATE users
				SET technology_stack = JSONB_SET(technology_stack, '{techStack}', :tech_stack),
					soft_skills      = JSONB_SET(soft_skills, '{softSkills}', :soft_skills),
					hobbies          = JSONB_SET(hobbies, '{hobbies}', :hobbies)
				WHERE users.id = 1
				RETURNING*;`,
		request,
	)

	if err != nil {
		return
	}
	return
}
