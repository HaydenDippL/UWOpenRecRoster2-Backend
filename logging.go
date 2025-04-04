package main

import (
    "UWOpenRecRoster2-Backend/models"
    "time"
    "fmt"
    "gorm.io/gorm"
    "errors"
    "github.com/google/uuid"
)

func log_event(userId string, sessionId string, date time.Time) (string, string, error) {
    // validate/get/create userId
    user, userErr := getUser(userId)
    createdNewUser := false
    if errors.Is(userErr, gorm.ErrRecordNotFound) {
        user, userErr = createUser()
        createdNewUser = true
        if userErr != nil {
            return "", "", fmt.Errorf("error creating new user")
        }
    } else if userErr != nil {
        return "", "", fmt.Errorf("error retrieving user")
    }

    // validate/get/create sessionId
    var session = models.Session{}
    var sessionErr error = nil
    if !createdNewUser {
        session, sessionErr = getSession(sessionId)
    }

    if createdNewUser || errors.Is(sessionErr, gorm.ErrRecordNotFound) {
        session, sessionErr = createSession(user.UserID)
        if sessionErr != nil {
            return "", "", fmt.Errorf("error creating new session")
        }
    } else if sessionErr != nil {
        return "", "", fmt.Errorf("error retrieving session")
    }
    
    // create query
    query := models.Query{
        SessionID: session.SessionID,
        ScheduleDate: date,
        QueriedTime: time.Now(),
    }

    result := DB.Create(&query)
    if result.Error != nil {
        return "", "", result.Error
    }

    return user.UserID, session.SessionID, nil
}

func getUser(userId string) (models.User, error) {
    if userId == "" {
        return models.User{}, nil
    }

    var user models.User
    result := DB.Where("user_id = ?", userId).First(&user)
    if result.Error != nil {
        return models.User{}, result.Error
    }

    return user, nil
}

func createUser() (models.User, error) {
    var user models.User
    for i := 0; i < 10; i++ {
        uuid := uuid.New().String()

        user = models.User{
            UserID: uuid,
            Created: time.Now(),
        }

        result := DB.Create(&user)
        
        if result.Error != nil && !errors.Is(result.Error, gorm.ErrDuplicatedKey) {
            return models.User{}, result.Error
        }
    }

    return user, nil
}

func getSession(sessionId string) (models.Session, error) {
    if sessionId == "" {
        return models.Session{}, nil
    }

    var session models.Session
    result := DB.Where("session_id = ?", sessionId).First(&session)
    if result.Error != nil {
        return models.Session{}, result.Error
    }

    return session, nil
}

func createSession(userId string) (models.Session, error) {
    var session models.Session
    for i := 0; i < 10; i++ {
        uuid := uuid.New().String()

        session = models.Session{
            UserID: userId,
            SessionID: uuid,
            Created: time.Now(),
        }

        result := DB.Create(&session)

        if result.Error != nil && !errors.Is(result.Error, gorm.ErrDuplicatedKey) {
            return models.Session{}, result.Error
        }
    }

    return session, nil
}
