package models

import (
	"time"
)

type User struct {
	UserID  string    `gorm:"type:uuid;unique;not null"`
	Created time.Time `gorm:"type:timestamp;not null"`
}

type Session struct {
	User      User      `gorm:"foreignKey:UserID;references:UserID"`
	SessionID string    `gorm:"type:uuid;unique;not null"`
	Created   time.Time `gorm:"type:timestamp;not null"`
}

type Query struct {
	Sessions Session   `gorm:"foreignKey:SessionID;references:SessionID"`
	Gym      string    `gorm:"type:varchar(100);not null;check:Gym IN ('nick', 'bakke')"`
	Facility string    `gorm:"type:varchar(100);not null"`
	Created  time.Time `gorm:"type:timestamp;not null"`
}
