package userinfo

import (
	"strconv"
	"strings"
	"time"
)

type UserInfo struct {
	UploadUsed    int64
	DownloadUsed  int64
	TotalTransfer int64
	ExpireAt      *time.Time
}

func ParseSubscriptionUserinfo(header string) *UserInfo {
	if header == "" {
		return nil
	}

	info := &UserInfo{}
	parts := strings.Split(header, ";")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "upload":
			if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				info.UploadUsed = v
			}
		case "download":
			if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				info.DownloadUsed = v
			}
		case "total":
			if v, err := strconv.ParseInt(value, 10, 64); err == nil {
				info.TotalTransfer = v
			}
		case "expire":
			if v, err := strconv.ParseInt(value, 10, 64); err == nil && v > 0 {
				t := time.Unix(v, 0)
				info.ExpireAt = &t
			}
		}
	}

	if info.UploadUsed == 0 && info.DownloadUsed == 0 && info.TotalTransfer == 0 && info.ExpireAt == nil {
		return nil
	}

	return info
}
