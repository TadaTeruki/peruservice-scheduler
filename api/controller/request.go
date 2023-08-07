package controller

import "time"

type LoginRequest struct {
	AdminID  string `json:"adminID"`
	Password string `json:"password"`
}

type PostScheduleRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Allday      bool      `json:"allday"`
	Tags        []string  `json:"tags"`
	Properties  []string  `json:"properties"`
	IsPublic    bool      `json:"isPublic"`
}

type PutScheduleRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Allday      bool      `json:"allday"`
	Tags        []string  `json:"tags"`
	Properties  []string  `json:"properties"`
	IsPublic    bool      `json:"isPublic"`
}
