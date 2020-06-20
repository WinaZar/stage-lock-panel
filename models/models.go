package models

import "github.com/jinzhu/gorm"

// Stage presents stage server model
type Stage struct {
	gorm.Model
	Name     string               `gorm:"type:varchar(80);unique;not null" json:"name" valid:"type(string),required,alphanum,length(1|15)"`
	LockCode string               `gorm:"type:varchar(80)" json:"code"`
	Locked   bool                 `gorm:"default:false;not null" json:"locked"`
	LockedBy string               `gorm:"type:varchar(80)" json:"locked_by"`
	Comment  string               `gorm:"type:varchar(500)" json:"comment"`
	History  []StageHistoryRecord `gorm:"foreignkey:Stage;association_foreignkey:Name"`
}

// StageHistoryRecord presents history of locks/unlocks
type StageHistoryRecord struct {
	gorm.Model
	Action   string `gorm:"not null" json:"action"`
	LockedBy string `json:"locked_by"`
	Comment  string `gorm:"type:varchar(500)" json:"comment"`
	Stage    string `gorm:"not null" json:"-"`
}
