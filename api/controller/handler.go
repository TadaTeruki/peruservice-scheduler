package controller

import (
	"net/http"
	"time"

	"github.com/TadaTeruki/peruservice-scheduler/config"
	"github.com/TadaTeruki/peruservice-scheduler/pkg/timeconv"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/lib/pq"
)

type Handler struct {
	db     *sqlx.DB
	config *config.ServerConfig
}

func NewHandler(conf *config.ServerConfig) (*Handler, error) {
	db, err := sqlx.Open(
		"postgres",
		"host="+conf.EnvConf.DBHost+
			" port="+conf.EnvConf.DBPort+
			" user="+conf.EnvConf.DBUser+
			" password="+conf.EnvConf.DBPassWord+
			" dbname="+conf.EnvConf.DBName+
			" sslmode=disable",
	)
	if err != nil {
		return nil, err
	}
	return &Handler{
		db:     db,
		config: conf,
	}, nil
}

func (h *Handler) PostSchedule(c echo.Context) error {
	permission := c.Get(string(ContextKeyPermission)).(Permission)
	if permission != PermissionAdmin {
		return c.JSON(http.StatusForbidden, ErrorResponse{
			Status:  http.StatusForbidden,
			Message: "permission denied",
		})
	}

	var request PostScheduleRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "failed to bind request",
		})
	}

	scheduleID := uuid.New().String()

	schedule := Schedule{
		ID:          scheduleID,
		Title:       request.Title,
		Description: request.Description,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		Allday:      request.Allday,
		Tags:        pq.StringArray(request.Tags),
		Properties:  pq.StringArray(request.Properties),
		IsPublic:    request.IsPublic,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := h.db.NamedExec(
		"INSERT INTO schedules (id, title, description, start_date, end_date, allday, tags, properties, is_public, created_at, updated_at) VALUES (:id, :title, :description, :start_date, :end_date, :allday, :tags, :properties, :is_public, :created_at, :updated_at)",
		schedule,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to insert schedule: " + err.Error(),
		})
	}

	response := PostScheduleResponse{ScheduleID: scheduleID}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetSchedule(c echo.Context) error {
	permission := c.Get(string(ContextKeyPermission)).(Permission)
	scheduleID := c.Param("scheduleID")

	schedule := Schedule{}
	if err := h.db.Get(&schedule, "SELECT * FROM schedules WHERE id = $1", scheduleID); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to get schedule from db",
		})
	}

	if !schedule.IsPublic && permission != PermissionAdmin {
		schedule.SetDetailsHidden()
	}

	response := GetScheduleResponse{
		ScheduleID:  schedule.ID,
		Title:       schedule.Title,
		Description: schedule.Description,
		StartDate:   schedule.StartDate,
		EndDate:     schedule.EndDate,
		Allday:      schedule.Allday,
		Tags:        schedule.Tags,
		Properties:  schedule.Properties,
		IsPublic:    schedule.IsPublic,
		CreatedAt:   schedule.CreatedAt,
		UpdatedAt:   schedule.UpdatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetScheduleList(c echo.Context) error {
	permission := c.Get(string(ContextKeyPermission)).(Permission)

	startDate, err := timeconv.StringToTime(c.QueryParam("startDate"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "failed to parse start date",
		})
	}

	endDate, err := timeconv.StringToTime(c.QueryParam("endDate"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "failed to parse end date",
		})
	}

	_, _ = startDate, endDate

	scheduleList := []Schedule{}
	if err := h.db.Select(&scheduleList, "SELECT * FROM schedules WHERE end_date >= $1 AND start_date <= $2", startDate, endDate); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to get schedule list from db",
		})
	}

	scheduleListResponse := []GetScheduleResponse{}
	for _, schedule := range scheduleList {
		if !schedule.IsPublic && permission != PermissionAdmin {
			schedule.SetDetailsHidden()
		}
		scheduleListResponse = append(scheduleListResponse, GetScheduleResponse{
			ScheduleID:  schedule.ID,
			Title:       schedule.Title,
			Description: schedule.Description,
			StartDate:   schedule.StartDate,
			EndDate:     schedule.EndDate,
			Allday:      schedule.Allday,
			Tags:        schedule.Tags,
			Properties:  schedule.Properties,
			IsPublic:    schedule.IsPublic,
			CreatedAt:   schedule.CreatedAt,
			UpdatedAt:   schedule.UpdatedAt,
		})
	}

	response := GetScheduleListResponse{
		ScheduleList: scheduleListResponse,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) PutSchedule(c echo.Context) error {
	permission := c.Get(string(ContextKeyPermission)).(Permission)
	if permission != PermissionAdmin {
		return c.JSON(http.StatusForbidden, ErrorResponse{
			Status:  http.StatusForbidden,
			Message: "permission denied",
		})
	}

	scheduleID := c.Param("scheduleID")

	var request PutScheduleRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	schedule := Schedule{
		ID:          scheduleID,
		Title:       request.Title,
		Description: request.Description,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		Allday:      request.Allday,
		Tags:        pq.StringArray(request.Tags),
		Properties:  pq.StringArray(request.Properties),
		IsPublic:    request.IsPublic,
		UpdatedAt:   time.Now(),
	}

	_, err := h.db.NamedExec(
		"UPDATE schedules SET title = :title, description = :description, start_date = :start_date, end_date = :end_date, allday = :allday, tags = :tags, properties = :properties, is_public = :is_public, updated_at = :updated_at WHERE id = :id",
		schedule,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to update schedule: " + err.Error(),
		})
	}
	response := PutScheduleResponse{ScheduleID: scheduleID}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteSchedule(c echo.Context) error {
	permission := c.Get(string(ContextKeyPermission)).(Permission)
	if permission != PermissionAdmin {
		return c.JSON(http.StatusForbidden, ErrorResponse{
			Status:  http.StatusForbidden,
			Message: "permission denied",
		})
	}

	scheduleID := c.Param("scheduleID")

	_, err := h.db.Exec("DELETE FROM schedules WHERE id = $1", scheduleID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete schedule: " + err.Error(),
		})
	}

	response := PutScheduleResponse{ScheduleID: scheduleID}

	return c.JSON(http.StatusOK, response)
}
