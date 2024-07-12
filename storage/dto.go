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
	USERID   uint   `gorm:"primary_key" json:"user_id"`
	EMAIL    string `gorm:"unique" json:"email"`
	LOGIN    string `gorm:"unique" json:"login"`
	USERNAME string `json:"username"`
	SURNAME  string `json:"surname"`
	PASSWORD string `json:"password" validate:"required"`
	ISADMIN  bool   `gorm:"default:false" json:"isadmin"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateRequest struct {
	Username string `json:"username"`
	Surname  string `json:"surname"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required" `
	Login    string `json:"login" `
}
