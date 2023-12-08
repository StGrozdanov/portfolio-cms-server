package auth

type UserAuthData struct {
	Username string `db:"username" json:"username" valid:"required,minstringlength(3)"`
	Password string `db:"password" json:"password" valid:"required,minstringlength(5)"`
}
