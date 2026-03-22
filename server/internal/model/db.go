package model

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	if err := autoMigrate(); err != nil {
		return err
	}
	return migrateData()
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&User{},
		&Subscription{},
		&Config{},
		&Rule{},
		&Script{},
	)
}

func migrateData() error {
	if err := migrateSubscriptions(); err != nil {
		return err
	}
	return migrateRulesAndScripts()
}

func migrateSubscriptions() error {
	var count int64
	DB.Model(&Subscription{}).Where("source_type = ? OR source_type = ''", "").Count(&count)
	if count == 0 {
		return nil
	}
	return DB.Model(&Subscription{}).Where("source_type = ? OR source_type = ''", "").Updates(map[string]interface{}{
		"source_type": SourceTypeRemote,
		"use_proxy":   false,
		"skip_cert":   false,
	}).Error
}

func migrateRulesAndScripts() error {
	var rules []Rule
	if err := DB.Where("subscription_id = 0 OR subscription_id IS NULL").Find(&rules).Error; err != nil {
		return err
	}
	if len(rules) == 0 {
		return nil
	}
	var activeSub Subscription
	if err := DB.Where("active = ?", true).First(&activeSub).Error; err != nil {
		var sub Subscription
		if err := DB.First(&sub).Error; err != nil {
			localSub := &Subscription{
				Name:       "默认配置",
				SourceType: SourceTypeLocal,
				Content:    "",
				Active:     false,
			}
			if err := DB.Create(localSub).Error; err != nil {
				return err
			}
			activeSub = *localSub
		} else {
			activeSub = sub
		}
	}
	for _, rule := range rules {
		rule.SubscriptionID = activeSub.ID
		if err := DB.Save(&rule).Error; err != nil {
			return err
		}
	}
	var scripts []Script
	if err := DB.Where("subscription_id = 0 OR subscription_id IS NULL").Find(&scripts).Error; err != nil {
		return err
	}
	for _, script := range scripts {
		script.SubscriptionID = activeSub.ID
		if err := DB.Save(&script).Error; err != nil {
			return err
		}
	}
	return nil
}

var _ = time.Time{}
