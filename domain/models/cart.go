package models

import (
	"time"
)

type Cart struct {
	UserId    int64     `json:"user_id" gorm:"column:user_id;not_null"`
	SKUId     int64     `json:"sku_id" gorm:"column:sku_id;not_null"`
	Count     int32     `json:"count" gorm:"column:count;not_null"`
	Timestamp time.Time `json:"timestamp" gorm:"column:timestamp;not_null"`
	Status    int8      `json:"status" gorm:"column:status;not_null;default:0"`
}