package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Data represents the nested data structure in the request
type Data struct {
	BuildingId       int    `json:"BuildingId"`
	Title            string `json:"Title"`
	Format           int    `json:"Format"`
	DropEventsInPast bool   `json:"DropEventsInPast"`
	EncryptD         string `json:"EncryptD"`
}

// RequestBody represents the complete request body structure
type RequestBody struct {
	Date string `json:"date"`
	Data Data   `json:"data"`
}

type Gym struct {
	title   string
	id      int
	encrypt string
}

var (
	bakke = Gym{
		title:   "Bakke Recreation and Wellbeing Center",
		id:      1112,
		encrypt: "https://uwmadison.emscloudservice.com/web/CustomBrowseEvents.aspx?data=meoZqrqZMvHKSLWaHS%2f4bjdroAMc1geNvtL12O1chw1fIP%2bOGy79Y1bkm2DPPKqmpSFHyPvFHX3LAJJHEfBPycyxctYlpcHD4rIwd%2byAtBNWXsKhJT9UDchzs%2bSc3Ze6JFHimlPlQrL2Jk7LFEkj3FoTWmA0BKzQQk0%2beDFO2IBZSiNnDXPGZQ%3d%3d",
	}
	nick = Gym{
		title:   "Nicholas Recreation Center",
		id:      1109,
		encrypt: "https://uwmadison.emscloudservice.com/web/CustomBrowseEvents.aspx?data=RtFXo1hK2Mh0UPlwkh3Aua7auJ66NvvBNBlUULUwM7vu4XjCwc5WoatHUWdz5pRofwluz9ZmHCNbHsgQ9uEDZjArIem0ShC%2fuM4gJbohNWkNGhzqKkAwrHDWzuEbcQxjHc8CzLweyL05oQ7ToCjKkM5TC%2b639V3qHwqgx1EhbWU%3d",
	}
)

const RECWELL_SCHEDULES_URL string = "https://uwmadison.emscloudservice.com/web/AnonymousServersApi.aspx/CustomBrowseEvents"

func fetchSchedule(gym Gym, date string) (string, error) {
	body := RequestBody{
		Date: date,
		Data: Data{
			BuildingId:       gym.id,
			Title:            gym.title,
			Format:           0,
			DropEventsInPast: false,
			EncryptD:         gym.encrypt,
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := http.Post(RECWELL_SCHEDULES_URL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	schedule, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println(string(schedule))

	return string(schedule), nil
}
