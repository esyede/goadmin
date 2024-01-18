package model

import (
	"time"

	"gorm.io/gorm"
)

type OperationLog struct {
	gorm.Model
	Username   string    `gorm:"type:varchar(20);comment:'Username'" json:"username"`
	Ip         string    `gorm:"type:varchar(20);comment:'IP address'" json:"ip"`
	IpLocation string    `gorm:"type:varchar(20);comment:'IP location'" json:"ipLocation"`
	Method     string    `gorm:"type:varchar(20);comment:'Request method'" json:"method"`
	Path       string    `gorm:"type:varchar(100);comment:'Access path'" json:"path"`
	Desc       string    `gorm:"type:varchar(100);comment:'Description'" json:"desc"`
	Status     int       `gorm:"type:int(4);comment:'Response status code'" json:"status"`
	StartTime  time.Time `gorm:"type:datetime(3);comment:'Start time'" json:"startTime"`
	TimeCost   int64     `gorm:"type:int(6);comment:'Request time (ms)'" json:"timeCost"`
	UserAgent  string    `gorm:"type:varchar(20);comment:'User agent'" json:"userAgent"`
}
