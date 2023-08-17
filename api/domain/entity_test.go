package domain

import (
	"testing"
	"time"
)

func TestConstantScheduleIntoSchedules(t *testing.T) {
	s := &Schedule{
		StartDate: time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2023, 1, 7, 0, 0, 0, 0, time.UTC),
		Constant: &Constant{
			StartDate:    time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC),
			EndDate:      time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
			IntervalDays: 5,
		},
	}

	result := s.ConstantScheduleIntoSchedules(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC))

	expected := []Schedule{
		{
			StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			StartDate: time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2023, 1, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			StartDate: time.Date(2023, 1, 11, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2023, 1, 12, 0, 0, 0, 0, time.UTC),
		},
	}

	if len(expected) != len(result) {
		t.Errorf("expected: %v, result: %v", expected, result)
	}
	for i, v := range expected {
		if v.StartDate != result[i].StartDate || v.EndDate != result[i].EndDate {
			t.Errorf("expected: %v, result: %v", expected, result)
		}
	}
}
