package model

import (
	"go-sales/pkg/util"
	"time"

	"gorm.io/gorm"
)

type CompanyGlobalContact struct {
	ID        *util.UUID `gorm:"column:id;type:uuid;primary_key"`
	CompanyID *util.UUID `gorm:"column:company_id;type:uuid"`
	Name      string     `gorm:"column:name;type:varchar(255)"`
	Email     *string    `gorm:"column:email;type:varchar(150)"`
	Phone     *string    `gorm:"column:phone;type:varchar(20)"`
	CGC       *string    `gorm:"column:cgc;type:varchar(40)"`
}

func (CompanyGlobalContact) TableName() string {
	return "company_global_contacts"
}

type CompanyGlobalAddress struct {
	ID               *util.UUID `gorm:"column:id;type:uuid;primary_key"`
	CompanyID        *util.UUID `gorm:"column:company_id;type:uuid"`
	Street           string     `gorm:"column:street;type:varchar(255)"`
	StreetNumber     *string    `gorm:"column:street_number;type:varchar(50)"`
	StreetComplement *string    `gorm:"column:street_complement;type:varchar(255)"`
	City             string     `gorm:"column:city;type:varchar(100)"`
	State            string     `gorm:"column:state;type:varchar(100)"`
	PostalCode       string     `gorm:"column:postal_code;type:varchar(20)"`
	Country          string     `gorm:"column:country;type:varchar(100)"`
}

func (CompanyGlobalAddress) TableName() string {
	return "company_global_addresses"
}

type CompanyGlobal struct {
	ID          *util.UUID `gorm:"column:id;type:uuid;primary_key"`
	Name        string     `gorm:"column:name;type:varchar(255)"`
	SocialName  string     `gorm:"column:social_name;type:varchar(255)"`
	Description *string    `gorm:"column:description;type:text"`
	CGC         string     `gorm:"column:cgc;type:varchar(40)"`
	Enabled     bool       `gorm:"column:enabled;type:boolean"`
	Email       *string    `gorm:"column:email;type:varchar(150)"`

	Address  *CompanyGlobalAddress   `gorm:"foreignKey:CompanyID;references:ID"`
	Contacts []*CompanyGlobalContact `gorm:"foreignKey:CompanyID"`

	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;autoCreateTime;<-:create"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz"`
}

func (CompanyGlobal) TableName() string {
	return "company_globals"
}
