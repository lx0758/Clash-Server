package model

import "time"

type SourceType string

const (
	SourceTypeRemote SourceType = "remote"
	SourceTypeLocal  SourceType = "local"
)

type Subscription struct {
	Model
	Name          string     `gorm:"size:255;not null" json:"name"`
	SourceType    SourceType `gorm:"size:16;not null;default:remote" json:"source_type"`
	URL           string     `gorm:"size:1024" json:"url"`
	Interval      int        `gorm:"default:60" json:"interval"`
	Content       string     `gorm:"type:text" json:"content"`
	Active        bool       `gorm:"default:false" json:"active"`
	UseProxy      bool       `gorm:"default:false" json:"use_proxy"`
	UserAgent     string     `gorm:"size:255" json:"user_agent"`
	SkipCert      bool       `gorm:"default:false" json:"skip_cert"`
	LastRefresh   *time.Time `json:"last_refresh"`
	NodeCount     int        `gorm:"default:0" json:"node_count"`
	UploadUsed    int64      `gorm:"default:0" json:"upload_used"`
	DownloadUsed  int64      `gorm:"default:0" json:"download_used"`
	TotalTransfer int64      `gorm:"default:0" json:"total_transfer"`
	ExpireAt      *time.Time `json:"expire_at"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
