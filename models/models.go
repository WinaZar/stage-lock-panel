package models

import (
	"time"
)

// Model is a base model
type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// Stage presents stage server model
type Stage struct {
	Model
	Name     string               `gorm:"type:varchar(80);unique;not null" json:"name" valid:"type(string),alphanum,length(3|15)"`
	LockCode string               `gorm:"type:varchar(80)" json:"code"`
	Locked   bool                 `gorm:"default:false;not null" json:"locked"`
	LockedBy string               `gorm:"type:varchar(80)" json:"locked_by"`
	Comment  string               `gorm:"type:varchar(500)" json:"comment"`
	History  []StageHistoryRecord `gorm:"foreignkey:Stage;references:Name"`
}

// StageHistoryRecord presents history of locks/unlocks
type StageHistoryRecord struct {
	Model
	Action   string `gorm:"not null" json:"action"`
	LockedBy string `json:"locked_by"`
	Comment  string `gorm:"type:varchar(500)" json:"comment"`
	Stage    string `gorm:"not null" json:"-"`
}
