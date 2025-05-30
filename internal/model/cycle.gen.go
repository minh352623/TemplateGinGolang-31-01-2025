// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameCycle = "cycle"

// Cycle mapped from table <cycle>
type Cycle struct {
	ID          int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Status      string    `gorm:"column:status;not null;default:draft" json:"status"`
	Sort        int32     `gorm:"column:sort" json:"sort"`
	UserCreated string    `gorm:"column:user_created" json:"user_created"`
	DateCreated time.Time `gorm:"column:date_created" json:"date_created"`
	UserUpdated string    `gorm:"column:user_updated" json:"user_updated"`
	DateUpdated time.Time `gorm:"column:date_updated" json:"date_updated"`
	Key         string    `gorm:"column:key" json:"key"`
	Value       int32     `gorm:"column:value" json:"value"`
}

// TableName Cycle's table name
func (*Cycle) TableName() string {
	return TableNameCycle
}
