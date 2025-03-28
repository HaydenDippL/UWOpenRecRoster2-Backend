package models

import (
	"time"
)

type User struct {
	UserID   string    `gorm:"type:uuid;unique;not null;primaryKey"`
	Created  time.Time `gorm:"type:timestamp;not null"`
	Sessions []Session `gorm:"foreignKey:UserID"` // Has many Sessions
}

type Session struct {
	SessionID string    `gorm:"type:uuid;unique;not null;primaryKey"`
	UserID    string    `gorm:"type:uuid;not null"` // Belongs to User
	Created   time.Time `gorm:"type:timestamp;not null"`
	Queries   []Query   `gorm:"foreignKey:SessionID"` // Has many Queries
}

type Query struct {
	Id           int       `gorm:"primaryKey;autoIncrement"`
	SessionID    string    `gorm:"type:uuid;not null"` // Belongs to Session
	Gym          string    `gorm:"type:varchar(100);not null;check:Gym IN ('nick', 'bakke')"`
	Facility     string    `gorm:"type:varchar(100);not null"`
	ScheduleDate time.Time `gorm:"type:timestamp;not null"`
	QueriedTime  time.Time `gorm:"type:timestamp;not null"`
}
