package model

import "go-sales/pkg/util"

type Role struct {
	ID   util.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Name string    `gorm:"column:name;type:varchar(255);unique;not null"`
}

func (Role) TableName() string {
	return "roles"
}
