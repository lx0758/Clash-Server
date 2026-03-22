package model

type User struct {
	Model
	Username string `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
}

func (User) TableName() string {
	return "users"
}
