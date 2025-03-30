package main

import (
    "UWOpenRecRoster2-Backend/models"
	"fmt"
	"time"
    "errors"
    "log"
    "gorm.io/gorm/clause"
)

// memoizes the schedule if it is in the range of [-3, 14] days from now
// also will delete all entries that are older than 3 days
func memoSchedule(schedule models.ScheduleJSON, date time.Time, gym string) error {
    log.Printf("Memoizing %s %v schedule\n", gym, date)

    if gym != "bakke" && gym != "nick" {
        return errors.New("gym must be either \"bakke\" or \"nick\"")
    }

    now := time.Now()
    twoWeeksAhead := now.AddDate(0, 0, 14)
    threeDaysAgo := now.AddDate(0, 0, -3)

    if err := DB.Where("schedule_date < ?", threeDaysAgo).Delete(&models.Schedule{}).Error; err != nil {
        return errors.New("error deleting all schedules older than three days")
    }

    if date.After(twoWeeksAhead) || date.Before(threeDaysAgo) {
        log.Printf("Aborting Memoize of %s %v, outside of memo window [-3 days, 14 days]\n", gym, date)
        return nil
    }

    scheduleDBModel := models.Schedule{
        ScheduleDate: date,
        Gym:          gym,
        Created:      time.Now(),
        Schedule:     schedule,
    }

    // Attempt to insert or update the record
    result := DB.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "schedule_date"}, {Name: "gym"}},
        DoUpdates: clause.AssignmentColumns([]string{"created", "schedule"}),
    }).Create(&scheduleDBModel)

    if result.Error != nil {
        return fmt.Errorf("failed to save schedule: %w", result.Error)
    }

    return nil
}

func getSchedule(date string, gym string) (models.ScheduleJSON, error) {
    var schedule models.Schedule

    parsedDate, err := time.Parse("2006-01-02", date)
    if err != nil {
        return models.ScheduleJSON{}, fmt.Errorf("invalid date format: %w", err)
    }
    
    err = DB.Where("schedule_date = ? AND gym = ?", parsedDate, gym).First(&schedule).Error
    if err != nil {
        return models.ScheduleJSON{}, fmt.Errorf("no occurrence of %s schedule for %s: %w", gym, date, err)
    }

    now := time.Now()
    oneHourAgo := now.Add(-1 * time.Hour)
    fmt.Printf("Created: %v, Now: %v\n", schedule.Created, now)
    if schedule.Created.Before(oneHourAgo) {
        return models.ScheduleJSON{}, fmt.Errorf("schedule is stale")
    }

    // get the schedule from the schedule DB object
    return schedule.Schedule, nil
}
