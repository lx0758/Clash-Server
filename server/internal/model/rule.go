package model

type RuleMode string

const (
	RuleModeInsert RuleMode = "insert"
	RuleModeAppend RuleMode = "append"
)

type Rule struct {
	Model
	SubscriptionID uint     `gorm:"not null;index" json:"subscription_id"`
	Name           string   `gorm:"size:255;not null" json:"name"`
	Type           string   `gorm:"size:64;not null" json:"type"`
	Payload        string   `gorm:"size:1024;not null" json:"payload"`
	Proxy          string   `gorm:"size:255;not null" json:"proxy"`
	Enabled        bool     `gorm:"default:true" json:"enabled"`
	Mode           RuleMode `gorm:"size:16;default:append" json:"mode"`
	Priority       int      `gorm:"default:0" json:"priority"`
}

func (Rule) TableName() string {
	return "rules"
}
