package apiserver

import (
	"avito-internship/internal/store"
	"avito-internship/internal/utils"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (s *APIServer) handleCreateExperiment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		name := request.URL.Query().Get("name")
		result, err := store.CreateExperiment(name)
		if err != nil {
			writeError(writer, err)
			return
		}
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

func (s *APIServer) handleGetAllExperiment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		allExperiment, err := store.GetAllExperiment()
		if err != nil {
			writeError(writer, err)
			return
		}
		data, err := json.Marshal(allExperiment)
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

func (s *APIServer) handleDeleteExperiment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		name := request.URL.Query().Get("name")
		err := store.DeleteExperiment(name)
		if err != nil {
			writeError(writer, err)
			return
		}
		_, err = io.WriteString(writer, "deleted: "+name)
		if err != nil {
			writeError(writer, err)
			return
		}
	}
}

func (s *APIServer) handleAddEmployeeToExperiment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.ParseInt(request.URL.Query().Get("employeeId"), 10, 8)
		if err != nil {
			writeError(writer, err)
			return
		}
		a := request.URL.Query().Get("experiments")
		err = store.AddEmployeeToExperiment(strings.Split(a, ","), int(id))
		if err != nil {
			writeError(writer, err)
			return
		}
	}
}

func (s *APIServer) handleRemoveEmployeeFromExperiment() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.ParseInt(request.URL.Query().Get("employeeId"), 10, 8)
		if err != nil {
			writeError(writer, err)
			return
		}
		a := request.URL.Query().Get("experiments")
		err = store.RemoveEmployeeFromExperiment(strings.Split(a, ","), int(id))
		if err != nil {
			writeError(writer, err)
			return
		}
	}
}

func (s *APIServer) getReportById() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.ParseInt(request.URL.Query().Get("employeeId"), 10, 8)
		if err != nil {
			writeError(writer, err)
			return
		}
		result, err := store.GetReportById(int(id))
		if err != nil {
			writeError(writer, err)
			return
		}
		fileName, err := utils.GetCsvFile(result)
		if err != nil {
			writeError(writer, err)
			return
		}
		file, err := os.Open(fileName)
		if err != nil {
			writeError(writer, err)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				writeError(writer, err)
			}
			err = os.Remove(fileName)
			if err != nil {
				writeError(writer, err)
				return
			}
		}(file)
		writer.Header().Set("Content-Type", "text/csv")
		writer.Header().Set("Content-Disposition", "attachment; filename=file.csv")
		_, err = io.Copy(writer, file)
		if err != nil {
			writeError(writer, err)
			return
		}
	}
}
