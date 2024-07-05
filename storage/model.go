package storage

type TOKENS struct {
	IDTOKENS     uint `gorm:"primary_key"`
	USERID       uint `gorm:"unique"`
	ACCESSTOCKEN string
	REFRESHTOKEN string
	EXP          int64
	TIMECREATE   int64 `gorm:"autoCreateTime"`
}

type USERS struct {
	USERID   uint   `gorm:"primary_key"`
	EMAIL    string `gorm:"unique" json:"email"`
	LOGIN    string `gorm:"unique" json:"login"`
	USERNAME string `json:"username"`
	SURNAME  string `json:"surname"`
	PASSWORD string `json:"password" validate:"required"`
}
