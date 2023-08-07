package controller

import (
	"time"

	"github.com/lib/pq"
)

type Schedule struct {
	ID          string         `db:"id"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
	StartDate   time.Time      `db:"start_date"`
	EndDate     time.Time      `db:"end_date"`
	Allday      bool           `db:"allday"`
	Tags        pq.StringArray `db:"tags"`
	Properties  pq.StringArray `db:"properties"`
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
