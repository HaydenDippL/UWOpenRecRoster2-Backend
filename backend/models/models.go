package models

import (
	"time"
    "fmt"
    "database/sql/driver"
    "encoding/json"
)

// Database models

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
	ScheduleDate time.Time `gorm:"type:timestamp;not null"`
	QueriedTime  time.Time `gorm:"type:timestamp;not null"`
}

type Schedule struct {
	ScheduleDate time.Time    `gorm:"type:date;not null;primaryKey"`
    Created      time.Time `gorm:"type:timestamptz;not null"`
	Schedule     ScheduleResp  `gorm:"type:jsonb;not null"`
}

type ScheduleResp struct {
    Bakke FacilityEvents `json:"bakke"`
    Nick FacilityEvents `json:"nick"`
}

type FacilityEvents struct {
    Courts []Event `json:"courts"`
    Pool []Event `json:"pool"`
    Esports []Event `json:"esports"`
    MtMendota []Event `json:"mount_mendota"`
    IceRink []Event `json:"ice_rink"`
}

type Event struct {
    Name string `json:"name"`
    Location string `json:"location"`
    Start string `json:"start"`
    End string `json:"end"`
}

// Scan implements the sql.Scanner interface for ScheduleJSON
func (s *ScheduleResp) Scan(value interface{}) error {
    if value == nil {
        return nil
    }
    
    var data []byte
    switch v := value.(type) {
    case string:
        data = []byte(v)
    case []byte:
        data = v
    default:
        return fmt.Errorf("unsupported type: %T", value)
    }
    
    return json.Unmarshal(data, s)
}

// Value implements the driver.Valuer interface for ScheduleJSON
func (s ScheduleResp) Value() (driver.Value, error) {
    return json.Marshal(s)
}

