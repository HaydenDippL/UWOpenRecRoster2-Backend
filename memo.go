package main

import (
	"fmt"
	"time"
)

type Schedule struct {
	ScheduleDate string    `gorm:"type:date;not null"`
	Gym          string    `gorm:"type:varchar(100);not null;check:Gym IN ('nick', 'bakke')"`
	Created      time.Time `gorm:"type:timestamp;not null"`
	Events       Events    `gorm:"type:jsonb;not null"`
}

func memoSchedule(date string, gym string) error {

	return nil
}

func getSchedule(date string, gym string) (Events, error) {
	var schedule Schedule
	err := DB.Where("schedule_date = ? AND gym = ?", date, gym).First(&schedule).Error
	if err != nil {
		return Events{}, fmt.Errorf("No occurence of %s schedule for %s: %w", gym, date, err)
	}

	return schedule.Events, nil
}
