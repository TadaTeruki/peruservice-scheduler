package domain

import "time"

type IScheduleRepository interface {
	PostSchedule(schedule *Schedule) error
	GetSchedule(scheduleID string) (*Schedule, error)
	GetScheduleList(startDate time.Time, endDate time.Time) (*[]Schedule, error)
	PutSchedule(schedule *Schedule) error
	DeleteSchedule(scheduleID string) error
}
