package internal

import "encoding/json"

type BasicUserInfo struct {
	Email    string          `db:"email" json:"email"`
	CvLink   string          `db:"cv_link" json:"cvLink"`
	AboutMe  string          `db:"about_me" json:"aboutMe"`
	Partners json.RawMessage `db:"partners" json:"partners"`
	Carousel json.RawMessage `db:"carousel" json:"carousel"`
}

type UserSkills struct {
	TechStack  json.RawMessage `db:"tech_stack" json:"techStack"`
	SoftSkills json.RawMessage `db:"soft_skills" json:"softSkills"`
	Hobbies    json.RawMessage `db:"hobbies" json:"hobbies"`
}
