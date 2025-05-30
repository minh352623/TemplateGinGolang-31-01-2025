// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNamePlatformInterestRate = "platform_interest_rates"

// PlatformInterestRate mapped from table <platform_interest_rates>
type PlatformInterestRate struct {
	ID               int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Status           string    `gorm:"column:status;not null;default:draft" json:"status"`
	Sort             int32     `gorm:"column:sort" json:"sort"`
	UserCreated      string    `gorm:"column:user_created" json:"user_created"`
	DateCreated      time.Time `gorm:"column:date_created" json:"date_created"`
	UserUpdated      string    `gorm:"column:user_updated" json:"user_updated"`
	DateUpdated      time.Time `gorm:"column:date_updated" json:"date_updated"`
	PercentForAdmin  float64   `gorm:"column:percent_for_admin" json:"percent_for_admin"`
	LockTime         int32     `gorm:"column:lock_time" json:"lock_time"`
	PercentPrincipal float64   `gorm:"column:percent_principal" json:"percent_principal"`
	PercentInterest  float64   `gorm:"column:percent_interest" json:"percent_interest"`
	WalletInterest   int32     `gorm:"column:wallet_interest" json:"wallet_interest"`
	Name             string    `gorm:"column:name" json:"name"`
}

// TableName PlatformInterestRate's table name
func (*PlatformInterestRate) TableName() string {
	return TableNamePlatformInterestRate
}
