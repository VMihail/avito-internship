package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	if err := s.configLogger(); err != nil {
		return err
	}
	s.configRouter()
	s.logger.Info("Starting API server")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configRouter() {
	s.router.HandleFunc("/getAllEmployee", s.handleGetAllEmployee())
	s.router.HandleFunc("/createEmployee", s.handleCreateEmployee())
	s.router.HandleFunc("/getEmployeeById", s.handleGetEmployeeById())

	s.router.HandleFunc("/createExperiment", s.handleCreateExperiment())
	s.router.HandleFunc("/getAllExperiment", s.handleGetAllExperiment())
	s.router.HandleFunc("/deleteExperiment", s.handleDeleteExperiment())

	s.router.HandleFunc("/addEmployeeToExperiment", s.handleAddEmployeeToExperiment())
	s.router.HandleFunc("/removeEmployeeFromExperiment", s.handleRemoveEmployeeFromExperiment())
	s.router.HandleFunc("/getExperimentsById", s.handleGetExperimentsById())

	s.router.HandleFunc("/getReportById", s.getReportById())
	s.router.HandleFunc("/getPercentageOfEmployees", s.handleGetPercentageOfEmployees())
	s.router.HandleFunc("/addPercentageOfEmployeeToExperiment", s.handleAddPercentageOfEmployeeToExperiment())
	s.router.HandleFunc("/handleAddEmployeeToExperimentWithTTL", s.handleAddEmployeeToExperimentWithTTL())
}

func writeError(writer http.ResponseWriter, err error) {
	_, _ = io.WriteString(writer, err.Error())
}
