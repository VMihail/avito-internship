package apiserver

import (
	"avito-internship/internal/store"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func (s *APIServer) handleCreateEmployee() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		employee, err := store.CreateEmployee()
		if err != nil {
			writeError(writer, err)
			return
		}
		data, err := json.Marshal(employee)
		if err != nil {
			writeError(writer, err)
			return
		}
		_, err = io.WriteString(writer, string(data))
		if err != nil {
			writeError(writer, err)
			return
		}
	}
}

func (s *APIServer) handleGetAllEmployee() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		allEmployees, err := store.GetAllEmployee()
		if err != nil {
			writeError(writer, err)
			return
		}
		data, err := json.Marshal(allEmployees)
		if err != nil {
			writeError(writer, err)
			return
		}
		_, err = io.WriteString(writer, string(data))
		if err != nil {
			writeError(writer, err)
			return
		}
	}
}

func (s *APIServer) handleGetEmployeeById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 8)
		if err != nil {
			writeError(writer, err)
			return
		}
		result, err := store.GetById(int(id))
		data, err := json.Marshal(result)
		if err != nil {
			writeError(writer, err)
			return
		}
		_, err = io.WriteString(writer, string(data))
		if err != nil {
			writeError(writer, err)
			return
		}
	}
}

func (s *APIServer) handleGetExperimentsById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.ParseInt(request.URL.Query().Get("id"), 10, 8)
		if err != nil {
			writeError(writer, err)
			return
		}
		result, err := store.GetExperimentsById(int(id))
		data, err := json.Marshal(result)
		if err != nil {
			writeError(writer, err)
			return
		}
		_, err = io.WriteString(writer, string(data))
		if err != nil {
			writeError(writer, err)
			return
		}
	}
}
