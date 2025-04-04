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
func memoSchedule(schedule models.ScheduleResp, date time.Time) error {
    log.Printf("Memoizing %s %v schedule\n", date)

    now := time.Now()
    twoWeeksAhead := now.AddDate(0, 0, 14)
    threeDaysAgo := now.AddDate(0, 0, -3)

    if err := DB.Where("schedule_date < ?", threeDaysAgo).Delete(&models.Schedule{}).Error; err != nil {
        return errors.New("error deleting all schedules older than three days")
    }

    if date.After(twoWeeksAhead) || date.Before(threeDaysAgo) {
        log.Printf("Aborting Memoize of %s %v, outside of memo window [-3 days, 14 days]\n", date)
        return nil
    }

    scheduleDBModel := models.Schedule{
        ScheduleDate: date,
        Created:      time.Now(),
        Schedule:     schedule,
    }

    // Attempt to insert or update the record
    result := DB.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "schedule_date"}},
        DoUpdates: clause.AssignmentColumns([]string{"created", "schedule"}),
    }).Create(&scheduleDBModel)

    if result.Error != nil {
        return fmt.Errorf("failed to save schedule: %w", result.Error)
    }

    return nil
}

func getSchedule(date string) (models.ScheduleResp, error) {
    var schedule models.Schedule

    parsedDate, err := time.Parse("2006-01-02", date)
    if err != nil {
        return models.ScheduleResp{}, fmt.Errorf("invalid date format: %w", err)
    }
    
    err = DB.Where("schedule_date = ?", parsedDate).First(&schedule).Error
    if err != nil {
        return models.ScheduleResp{}, fmt.Errorf("no occurrence of schedule for %s: %w", date, err)
    }

    now := time.Now()
    oneHourAgo := now.Add(-1 * time.Hour)
    fmt.Printf("Created: %v, Now: %v\n", schedule.Created, now)
    if schedule.Created.Before(oneHourAgo) {
        return models.ScheduleResp{}, fmt.Errorf("schedule is stale")
    }

    // get the schedule from the schedule DB object
    return schedule.Schedule, nil
}
