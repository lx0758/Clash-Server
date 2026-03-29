package model

type SubscriptionCustomization struct {
	Model
	SubscriptionID   uint   `gorm:"uniqueIndex;not null" json:"subscription_id"`
	ProxyInsert      string `gorm:"type:text" json:"proxy_insert"`
	ProxyAppend      string `gorm:"type:text" json:"proxy_append"`
	ProxyRemove      string `gorm:"type:text" json:"proxy_remove"`
	ProxyGroupInsert string `gorm:"type:text" json:"proxy_group_insert"`
	ProxyGroupAppend string `gorm:"type:text" json:"proxy_group_append"`
	ProxyGroupRemove string `gorm:"type:text" json:"proxy_group_remove"`
	RuleInsert       string `gorm:"type:text" json:"rule_insert"`
	RuleAppend       string `gorm:"type:text" json:"rule_append"`
	RuleRemove       string `gorm:"type:text" json:"rule_remove"`
	GlobalOverride   string `gorm:"type:text" json:"global_override"`
	Script           string `gorm:"type:text" json:"script"`
}

func (SubscriptionCustomization) TableName() string {
	return "subscription_customizations"
}
