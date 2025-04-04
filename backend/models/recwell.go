package models

// Request Format of the Recwell API call

type RequestBody struct {
    Date string `json:"date"`
	Data RequestData   `json:"data"`
}

type RequestData struct {
	BuildingId       int    `json:"BuildingId"`
	Title            string `json:"Title"`
	Format           int    `json:"Format"`
	DropEventsInPast bool   `json:"DropEventsInPast"`
	EncryptD         string `json:"EncryptD"`
}

// Response Format of the Recwell API call

type ResponseBody struct {
    Data string `json:"d"`
}

// Inner Response Format of the Recwell API call

type EventsRaw struct {
	Events []EventRaw `json:"DailyBookingResults"`
}

type EventRaw struct {
	EventName  string `json:"EventName"`
	Location   string `json:"Room"`
	EventStart string `json:"GmtStart"`
	EventEnd   string `json:"GmtEnd"`
}
