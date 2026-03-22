package model

type Script struct {
	Model
	SubscriptionID uint   `gorm:"not null;index" json:"subscription_id"`
	Name           string `gorm:"size:255;not null" json:"name"`
	Description    string `gorm:"size:512" json:"description"`
	Content        string `gorm:"type:text" json:"content"`
	Enabled        bool   `gorm:"default:true" json:"enabled"`
}

func (Script) TableName() string {
	return "scripts"
}
