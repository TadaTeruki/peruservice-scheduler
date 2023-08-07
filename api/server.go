package api

import (
	"log"
	"net/http"

	"github.com/TadaTeruki/peruservice-scheduler/api/controller"
	"github.com/TadaTeruki/peruservice-scheduler/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

type Server struct {
	Router  *echo.Echo
	config  *config.ServerConfig
	handler *controller.Handler
}

func NewServer(e *echo.Echo) *Server {
	conf, err := config.QueryServerConfig()
	if err != nil {
		log.Fatalf("failed to query server config: %v", err)
	}
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	handler, err := controller.NewHandler(conf)
	if err != nil {
		log.Fatalf("failed to create handler: %v", err)
	}
	return &Server{
		Router:  e,
		config:  conf,
		handler: handler,
	}
}

func (s *Server) Start() error {

	// Middleware
	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())

	// Allow origins
	var allow_origins []string
	if s.config.EnvConf.Mode == "PRODUCTION" {
		allow_origins = s.config.EnvConf.SchedulerAllowOrigins
	} else {
		allow_origins = []string{"*"}
	}

	// Set CORS
	s.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allow_origins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Authorization"},
		AllowCredentials: true,
	}))

	// Routes
	s.Router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "server: ok")
	})

	gschedule := s.Router.Group("/schedule")
	{
		gschedule.Use(controller.CheckPermissionMiddleware(s.config.EnvConf.PublicKeySrc))
		gschedule.GET("/:scheduleID", s.handler.GetSchedule)
		gschedule.POST("", s.handler.PostSchedule)
		gschedule.PUT("/:scheduleID", s.handler.PutSchedule)
		gschedule.DELETE("/:scheduleID", s.handler.DeleteSchedule)
	}
	glist := s.Router.Group("/schedules")
	{
		glist.Use(controller.CheckPermissionMiddleware(s.config.EnvConf.PublicKeySrc))
		glist.GET("", s.handler.GetScheduleList)
	}

	return s.Router.Start(":" + s.config.EnvConf.SchedulerPort)
}
