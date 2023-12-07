package internal

import "encoding/json"

type BasicUserInfo struct {
	Email    string          `db:"email" json:"email"`
	CvLink   string          `db:"cv_link" json:"cvLink"`
	AboutMe  string          `db:"about_me" json:"aboutMe"`
	Partners json.RawMessage `db:"partners" json:"partners"`
	Carousel json.RawMessage `db:"carousel" json:"carousel"`
}
