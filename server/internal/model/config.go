package model

type Config struct {
	Key   string `gorm:"primaryKey;size:64" json:"key"`
	Value string `gorm:"type:text" json:"value"`
}

func (Config) TableName() string {
	return "config"
}
