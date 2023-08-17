package infrastructure

import (
	"database/sql"
	"time"

	"github.com/TadaTeruki/peruservice-scheduler/api/domain"
	"github.com/TadaTeruki/peruservice-scheduler/config"
	"github.com/jmoiron/sqlx"
)

type ScheduleRepository struct {
	db *sqlx.DB
}

func NewScheduleRepository(config *config.ServerConfig) (*ScheduleRepository, error) {
	db, err := sqlx.Open(
		"postgres",
		"host="+config.EnvConf.DBHost+
			" port=5432"+
			" user="+config.EnvConf.DBUser+
			" password="+config.EnvConf.DBPassWord+
			" dbname="+config.EnvConf.DBName+
			" sslmode=disable",
	)
	if err != nil {
		return nil, err
	}
	return &ScheduleRepository{
		db: db,
	}, nil
}

func (r *ScheduleRepository) PostSchedule(schedule *domain.Schedule) error {
	tx, _ := r.db.Beginx()

	_, err := tx.NamedExec(
		"INSERT INTO schedules (id, title, description, start_date, end_date, tags, properties, is_public, created_at, updated_at) VALUES (:id, :title, :description, :start_date, :end_date, :tags, :properties, :is_public, :created_at, :updated_at)",
		schedule,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	if schedule.Constant != nil {
		_, err = tx.NamedExec(
			"INSERT INTO constants (schedule_id, start_date, end_date, interval_days) VALUES (:schedule_id, :start_date, :end_date, :interval_days)",
			schedule.Constant,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *ScheduleRepository) GetSchedule(scheduleID string) (*domain.Schedule, error) {
	var schedule domain.Schedule
	err := r.db.Get(&schedule, "SELECT * FROM schedules WHERE id = $1", scheduleID)
	if err != nil {
		return nil, err
	}

	var constant domain.Constant
	err = r.db.Get(&constant, "SELECT * FROM constants WHERE schedule_id = $1", scheduleID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		schedule.Constant = nil
	} else {
		schedule.Constant = &constant
	}

	return &schedule, nil
}

func (r *ScheduleRepository) GetScheduleList(startDate time.Time, endDate time.Time) (*[]domain.Schedule, error) {
	var scheduleList []domain.Schedule
	err := r.db.Select(&scheduleList, "SELECT * FROM schedules WHERE end_date >= $1 AND start_date <= $2 AND id NOT IN (SELECT schedule_id FROM constants)", startDate, endDate)
	if err != nil {
		return nil, err
	}

	var constantList []domain.Constant
	err = r.db.Select(&constantList, "SELECT * FROM constants")
	if err != nil {
		return nil, err
	}

	constantScheduleList := make([]domain.Schedule, len(constantList))
	for i, constant := range constantList {
		err := r.db.Get(&constantScheduleList[i], "SELECT * FROM schedules WHERE id = $1", constant.ScheduleID)
		if err != nil {
			return nil, err
		}
		constantScheduleList[i].Constant = &constantList[i]
	}

	for _, constantSchedule := range constantScheduleList {
		scs := constantSchedule.ConstantScheduleIntoSchedules(startDate, endDate)
		scheduleList = append(scheduleList, scs...)
	}

	return &scheduleList, nil
}

func (r *ScheduleRepository) PutSchedule(schedule *domain.Schedule) error {
	tx, _ := r.db.Beginx()
	_, err := tx.NamedExec(
		"UPDATE schedules SET title = :title, description = :description, start_date = :start_date, end_date = :end_date, tags = :tags, properties = :properties, is_public = :is_public, updated_at = :updated_at WHERE id = :id",
		schedule,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	if schedule.Constant != nil {
		_, err = tx.NamedExec(
			"INSERT INTO constants (schedule_id, start_date, end_date, interval_days) VALUES (:schedule_id, :start_date, :end_date, :interval_days) ON CONFLICT (schedule_id) DO UPDATE SET start_date = :start_date, end_date = :end_date, interval_days = :interval_days",
			schedule.Constant,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		_, err = tx.Exec("DELETE FROM constants WHERE schedule_id = $1", schedule.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *ScheduleRepository) DeleteSchedule(scheduleID string) error {
	tx, _ := r.db.Beginx()

	_, err := tx.Exec("DELETE FROM schedules WHERE id = $1", scheduleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM constants WHERE schedule_id = $1", scheduleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
