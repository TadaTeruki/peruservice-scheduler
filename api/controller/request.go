package controller

import "time"

type ConstantRequest struct {
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	IntervalDays int       `json:"intervalDays"`
}

type PostScheduleRequest struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	StartDate   time.Time        `json:"startDate"`
	EndDate     time.Time        `json:"endDate"`
	Tags        []string         `json:"tags"`
	Properties  []string         `json:"properties"`
	Constant    *ConstantRequest `json:"constant"`
	IsPublic    bool             `json:"isPublic"`
}

type PutScheduleRequest struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	StartDate   time.Time        `json:"startDate"`
	EndDate     time.Time        `json:"endDate"`
	Tags        []string         `json:"tags"`
	Properties  []string         `json:"properties"`
	Constant    *ConstantRequest `json:"constant"`
	IsPublic    bool             `json:"isPublic"`
}
