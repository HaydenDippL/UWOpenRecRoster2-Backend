package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
    "errors"
    "UWOpenRecRoster2-Backend/models"
)

type GymMetaData struct {
	title   string
	id      int
	encrypt string
}

var (
	bakke = GymMetaData{
		title:   "Bakke Recreation and Wellbeing Center",
		id:      1112,
		encrypt: "https://uwmadison.emscloudservice.com/web/CustomBrowseEvents.aspx?data=meoZqrqZMvHKSLWaHS%2f4bjdroAMc1geNvtL12O1chw1fIP%2bOGy79Y1bkm2DPPKqmpSFHyPvFHX3LAJJHEfBPycyxctYlpcHD4rIwd%2byAtBNWXsKhJT9UDchzs%2bSc3Ze6JFHimlPlQrL2Jk7LFEkj3FoTWmA0BKzQQk0%2beDFO2IBZSiNnDXPGZQ%3d%3d",
	}
	nick = GymMetaData{
		title:   "Nicholas Recreation Center",
		id:      1109,
		encrypt: "https://uwmadison.emscloudservice.com/web/CustomBrowseEvents.aspx?data=RtFXo1hK2Mh0UPlwkh3Aua7auJ66NvvBNBlUULUwM7vu4XjCwc5WoatHUWdz5pRofwluz9ZmHCNbHsgQ9uEDZjArIem0ShC%2fuM4gJbohNWkNGhzqKkAwrHDWzuEbcQxjHc8CzLweyL05oQ7ToCjKkM5TC%2b639V3qHwqgx1EhbWU%3d",
	}
)

const RECWELL_SCHEDULES_URL string = "https://uwmadison.emscloudservice.com/web/AnonymousServersApi.aspx/CustomBrowseEvents"

func fetchSchedules(date string) (models.ScheduleResp, error) {
    var schedule models.ScheduleResp

    for _, gym := range []string{"bakke", "nick"} {
        gymEvents, err := fetchSchedule(date, gym);
        if err != nil {
            return models.ScheduleResp{}, fmt.Errorf("error fetching the %s schedule", gym)
        }
        
        if gym == "bakke" {
            schedule.Bakke = gymEvents
        } else {
            schedule.Nick = gymEvents
        }
    }

    return schedule, nil
}

func fetchSchedule(date string, gym string) (models.FacilityEvents, error) {
    var gymMeta GymMetaData

    switch gym {
    case "bakke":
        gymMeta = bakke
    case "nick":
        gymMeta = nick
    default:
        return models.FacilityEvents{}, errors.New("gym must be either \"bakke\" or \"nick\"")
    }

	body := models.RequestBody{
		Date: date,
    	Data: models.RequestData{
			BuildingId:       gymMeta.id,
			Title:            gymMeta.title,
			Format:           0,
			DropEventsInPast: false,
			EncryptD:         gymMeta.encrypt,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return models.FacilityEvents{}, fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := http.Post(RECWELL_SCHEDULES_URL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return models.FacilityEvents{}, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.FacilityEvents{}, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	schedule, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.FacilityEvents{}, fmt.Errorf("failed to read response body: %w", err)
	}

	events, err := parseSchedule(schedule)
	if err != nil {
		return models.FacilityEvents{}, fmt.Errorf("failed to parse schedule: %w", err)
	}

    return events, nil
}

func parseSchedule(schedule []byte) (models.FacilityEvents, error) {
	var resp models.ResponseBody
	err := json.Unmarshal(schedule, &resp)
	if err != nil {
		return models.FacilityEvents{}, fmt.Errorf("error parsing JSON: %w", err)
	}

	var events models.EventsRaw
	err = json.Unmarshal([]byte(resp.Data), &events)
	if err != nil {
		return models.FacilityEvents{}, fmt.Errorf("error parsing JSON: %w", err)
	}

    var convertedEvents models.FacilityEvents = convertEventsToSchedule(events) 

	return convertedEvents, nil
}

const (
	court     = "court"
	mtMendota = "mount mendota"
	pool      = "pool"
	iceRink   = "ice rink"
	esports   = "esports"
)

func convertEventsToSchedule(events models.EventsRaw) models.FacilityEvents {
    schedule := models.FacilityEvents{}
    for _, eventRaw := range events.Events {
        location := strings.ToLower(strings.TrimSpace(eventRaw.Location))
        
        var event models.Event = transformAndDecodeRawEvent(eventRaw)
        
        if strings.Contains(location, court) {
            schedule.Courts = append(schedule.Courts, event)
        } else if strings.Contains(location, mtMendota) {
            schedule.MtMendota = append(schedule.MtMendota, event)
        } else if strings.Contains(location, pool) {
            schedule.Pool = append(schedule.Pool, event)
        } else if strings.Contains(location, iceRink) {
            schedule.IceRink = append(schedule.IceRink, event)
        } else if strings.Contains(location, esports) {
            schedule.Esports = append(schedule.Esports, event)   
        }
    }

    return schedule
}

func transformAndDecodeRawEvent(event models.EventRaw) models.Event {
    return models.Event{
        Name: html.UnescapeString(event.EventName),
        Location: html.UnescapeString(event.Location),
        Start: event.EventStart,
        End: event.EventEnd,
    }
}

