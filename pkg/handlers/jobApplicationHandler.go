package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/even44/JobsearchAPI/pkg/models"
	"github.com/even44/JobsearchAPI/pkg/stores"
	"github.com/gorilla/mux"
)

var JAH *JobApplicationsHandler

type JobApplicationsHandler struct {
	store  stores.JobApplicationStore
	logger *log.Logger
}

func NewJobApplicationHandler(s stores.JobApplicationStore) *JobApplicationsHandler {
	return &JobApplicationsHandler{
		store:  s,
		logger: log.New(os.Stdout, "[JOBAPPLICATION] ", log.Ldate+log.Ltime+log.Lmsgprefix),
	}
}

func (h JobApplicationsHandler) CreateJobApplication(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)

	var jobApplication models.JobApplication
	h.logger.Printf("Received request to create job application from: %s", r.Host)
	if err := json.NewDecoder(r.Body).Decode(&jobApplication); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	resultJobApplication, err := h.store.AddApplication(jobApplication)
	if err != nil {
		h.logger.Println(err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(resultJobApplication)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)

}
func (h JobApplicationsHandler) ListJobApplications(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)

	h.logger.Printf("Received request to list job applications from: %s", r.Host)
	jobapplications, err := h.store.ListApplications()
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting jobApplications list: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}
	jsonBytes, err := json.Marshal(jobapplications)
	if len(jobapplications) == 0 {
		jsonBytes = []byte("[]")
	}

	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}
func (h JobApplicationsHandler) GetJobApplication(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)
	strId := mux.Vars(r)["id"]
	h.logger.Printf("Received request to get job application with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	jobApplication, err := h.store.GetApplication(uint(id))

	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting jobApplication with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(jobApplication)

	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h JobApplicationsHandler) UpdateJobApplication(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)

	strId := mux.Vars(r)["id"]
	h.logger.Printf("Received request to update job application with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}
	oldJobApplication, err := h.store.GetApplication(uint(id))
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting jobApplication with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	var newJobApplication models.JobApplication

	if err := json.NewDecoder(r.Body).Decode(&newJobApplication); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	newJobApplication.ID = oldJobApplication.ID
	if newJobApplication.CompanyID == 0 {
		h.logger.Printf(
			"[UPDATE] Updated job application had companyid = 0, using old company id: %d",
			oldJobApplication.CompanyID)
		newJobApplication.CompanyID = oldJobApplication.CompanyID
	}

	err = h.store.UpdateApplication(uint(id), newJobApplication)

	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting jobApplication with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h JobApplicationsHandler) DeleteJobApplication(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)
	strId := mux.Vars(r)["id"]
	h.logger.Printf("Received request to delete job application with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	err = h.store.RemoveApplication(uint(id))

	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting jobApplication with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)

}
