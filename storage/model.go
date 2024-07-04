package storage

type TOCKENS struct {
	ID           uint `gorm:"primary_key"`
	USERID       uint
	ACCESSTOCKEN string
	REFRESHTOKEN string
	EXP          int64
}

type USERS struct {
	USERID   uint   `gorm:"primary_key"`
	USERNAME string `json:"username"`
	SURNAME  string `json:"surname"`
	EMAIL    string `json:"email"`
	LOGIN    string `json:"login"`
	PASSWORD string `json:"password" validate:"required"`
}
