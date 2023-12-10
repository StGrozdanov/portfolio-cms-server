package users

import "encoding/json"

type BasicUserInfo struct {
	Email    string          `db:"email" json:"email" valid:"required"`
	CvLink   string          `db:"cv_link" json:"cvLink" valid:"required"`
	AboutMe  string          `db:"about_me" json:"aboutMe" valid:"required"`
	Partners json.RawMessage `db:"partners" json:"partners" valid:"required"`
	Carousel json.RawMessage `db:"carousel" json:"carousel" valid:"required"`
}

type UserSkills struct {
	TechStack  json.RawMessage `db:"tech_stack" json:"techStack" valid:"required"`
	SoftSkills json.RawMessage `db:"soft_skills" json:"softSkills" valid:"required"`
	Hobbies    json.RawMessage `db:"hobbies" json:"hobbies" valid:"required"`
}

type JobsAndProjects struct {
	Jobs     json.RawMessage `db:"jobs" json:"jobs" valid:"required"`
	Projects json.RawMessage `db:"projects" json:"projects" valid:"required"`
}

type Socials struct {
	SocialMedia json.RawMessage `db:"social_media" json:"socialMedia" valid:"required"`
}
