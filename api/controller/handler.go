package controller

import (
	"net/http"
	"time"

	"github.com/TadaTeruki/peruservice-scheduler/api/domain"
	"github.com/TadaTeruki/peruservice-scheduler/pkg/timeconv"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/lib/pq"
)

type Handler struct {
	repository domain.IScheduleRepository
}

func NewHandler(repository domain.IScheduleRepository) (*Handler, error) {
	return &Handler{
		repository: repository,
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

	var constant *domain.Constant
	if request.Constant != nil {
		constant = &domain.Constant{
			ScheduleID:   scheduleID,
			StartDate:    request.Constant.StartDate,
			EndDate:      request.Constant.EndDate,
			IntervalDays: request.Constant.IntervalDays,
		}
	}

	schedule := domain.Schedule{
		ID:          scheduleID,
		Title:       request.Title,
		Description: request.Description,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		Tags:        pq.StringArray(request.Tags),
		Properties:  pq.StringArray(request.Properties),
		IsPublic:    request.IsPublic,
		Constant:    constant,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := h.repository.PostSchedule(&schedule)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to insert schedule",
		})
	}

	response := PostScheduleResponse{ScheduleID: scheduleID}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetSchedule(c echo.Context) error {
	permission := c.Get(string(ContextKeyPermission)).(Permission)
	scheduleID := c.Param("scheduleID")

	schedule, err := h.repository.GetSchedule(scheduleID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to get schedule from db",
		})
	}

	if !schedule.IsPublic && permission != PermissionAdmin {
		schedule.SetDetailsHidden()
	}

	var constant *ConstantResponse
	if schedule.Constant != nil {
		constant = &ConstantResponse{
			StartDate:    schedule.Constant.StartDate,
			EndDate:      schedule.Constant.EndDate,
			IntervalDays: schedule.Constant.IntervalDays,
		}
	}

	response := GetScheduleResponse{
		ScheduleID:  schedule.ID,
		Title:       schedule.Title,
		Description: schedule.Description,
		StartDate:   schedule.StartDate,
		EndDate:     schedule.EndDate,
		Tags:        schedule.Tags,
		Properties:  schedule.Properties,
		Constant:    constant,
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

	scheduleList, err := h.repository.GetScheduleList(startDate, endDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to get schedule list from db",
		})
	}

	scheduleListResponse := []GetScheduleResponse{}
	for _, schedule := range *scheduleList {
		if !schedule.IsPublic && permission != PermissionAdmin {
			schedule.SetDetailsHidden()
		}
		scheduleListResponse = append(scheduleListResponse, GetScheduleResponse{
			ScheduleID:  schedule.ID,
			Title:       schedule.Title,
			Description: schedule.Description,
			StartDate:   schedule.StartDate,
			EndDate:     schedule.EndDate,
			Tags:        schedule.Tags,
			Properties:  schedule.Properties,
			Constant:    nil,
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
		return c.JSON(http.StatusBadRequest, "invalid request")
	}

	var constant *domain.Constant
	if request.Constant != nil {
		constant = &domain.Constant{
			ScheduleID:   scheduleID,
			StartDate:    request.Constant.StartDate,
			EndDate:      request.Constant.EndDate,
			IntervalDays: request.Constant.IntervalDays,
		}
	}

	schedule := domain.Schedule{
		ID:          scheduleID,
		Title:       request.Title,
		Description: request.Description,
		StartDate:   request.StartDate,
		EndDate:     request.EndDate,
		Tags:        pq.StringArray(request.Tags),
		Properties:  pq.StringArray(request.Properties),
		Constant:    constant,
		IsPublic:    request.IsPublic,
		UpdatedAt:   time.Now(),
	}
	err := h.repository.PutSchedule(&schedule)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to update schedule",
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

	err := h.repository.DeleteSchedule(scheduleID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed to delete schedule",
		})
	}

	response := PutScheduleResponse{ScheduleID: scheduleID}

	return c.JSON(http.StatusOK, response)
}
