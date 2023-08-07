package controller

import "time"

type PostScheduleResponse struct {
	ScheduleID string `json:"scheduleID"`
}

type GetScheduleResponse struct {
	ScheduleID  string    `json:"scheduleID"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Allday      bool      `json:"allday"`
	Tags        []string  `json:"tags"`
	Properties  []string  `json:"properties"`
	IsPublic    bool      `json:"isPublic"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type GetScheduleListResponse struct {
	ScheduleList []GetScheduleResponse `json:"scheduleList"`
}

type PutScheduleResponse struct {
	ScheduleID string `json:"scheduleID"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
