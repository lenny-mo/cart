package models

import "time"

type Cart struct {
	ID        int64     `json:"id" gorm:"id;auto_increment;primary_key;not_null"`
	UserId    int64     `json:"user_id" gorm:"user_id;not_null"`
	SKUId     int64     `json:"sku_id" gorm:"sku_id;not_null"`
	Count     int32     `json:"count" gorm:"count;not_null"`
	Timestamp time.Time `json:"timestamp" gorm:"timestamp;not_null"`
}


