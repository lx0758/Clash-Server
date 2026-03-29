package model

import (
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
		&SubscriptionCustomization{},
	)
}

func migrateData() error {
	return migrateSubscriptions()
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
