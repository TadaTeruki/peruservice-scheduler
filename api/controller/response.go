package controller

import "time"

type ConstantResponse struct {
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	IntervalDays int       `json:"intervalDays"`
}

type PostScheduleResponse struct {
	ScheduleID string `json:"scheduleID"`
}

type GetScheduleResponse struct {
	ScheduleID  string            `json:"scheduleID"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	StartDate   time.Time         `json:"startDate"`
	EndDate     time.Time         `json:"endDate"`
	Tags        []string          `json:"tags"`
	Properties  []string          `json:"properties"`
	Constant    *ConstantResponse `json:"constant"`
	IsPublic    bool              `json:"isPublic"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
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
