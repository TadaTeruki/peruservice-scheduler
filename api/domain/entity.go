package domain

import (
	"time"

	"github.com/lib/pq"
)

type Constant struct {
	ScheduleID   string    `db:"schedule_id"`
	StartDate    time.Time `db:"start_date"`
	EndDate      time.Time `db:"end_date"`
	IntervalDays int       `db:"interval_days"`
}

type Schedule struct {
	ID          string         `db:"id"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
	StartDate   time.Time      `db:"start_date"`
	EndDate     time.Time      `db:"end_date"`
	Tags        pq.StringArray `db:"tags"`
	Properties  pq.StringArray `db:"properties"`
	Constant    *Constant      `db:"constant"`
	IsPublic    bool           `db:"is_public"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

func (s *Schedule) SetDetailsHidden() {
	s.Title = ""
	s.Description = ""
	s.Tags = pq.StringArray{}
	s.Properties = pq.StringArray{}
}

func modFloat(a float64, b float64) float64 {
	return a - b*float64(int(a/b))
}

func (s *Schedule) ConstantScheduleIntoSchedules(startDate time.Time, endDate time.Time) []Schedule {
	if endDate.Before(s.Constant.StartDate) || startDate.After(s.Constant.EndDate) {
		return []Schedule{}
	}

	repeatStartDate := s.Constant.StartDate
	if startDate.After(repeatStartDate) {
		repeatStartDate = startDate
	}
	repeatEndDate := s.Constant.EndDate
	if endDate.Before(repeatEndDate) {
		repeatEndDate = endDate
	}

	diff := s.StartDate.Sub(repeatStartDate)
	mod := modFloat(diff.Seconds(), float64(s.Constant.IntervalDays)*24.*60.*60.)
	firstStartDate := repeatStartDate.Add(time.Duration(mod) * time.Second)
	var schedules []Schedule
	for sdate := firstStartDate; sdate.Before(repeatEndDate); sdate = sdate.Add(time.Duration(s.Constant.IntervalDays) * 24 * time.Hour) {
		edate := sdate.Add(s.EndDate.Sub(s.StartDate))
		if edate.Before(repeatStartDate) {
			continue
		}
		schedules = append(schedules, Schedule{
			ID:          s.ID,
			Title:       s.Title,
			Description: s.Description,
			StartDate:   sdate,
			EndDate:     edate,
			Tags:        s.Tags,
			Properties:  s.Properties,
			Constant:    nil,
			IsPublic:    s.IsPublic,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		})
	}
	return schedules
}
