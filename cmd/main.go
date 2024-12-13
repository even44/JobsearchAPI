package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/even44/JobsearchAPI/pkg/jobApplications"
	"github.com/gorilla/mux"
)

func main() {

	const port int = 3001

	fmt.Printf("Jobsearch API running on port: %d\n", port)

	// Create the store and Jobapplication handler
	store := jobApplications.NewMariaDBStore()
	jobApplicationsHandler := NewJobApplicationHandler(store)

	// Create router
	router := mux.NewRouter()

	router.HandleFunc("/jobapplications", jobApplicationsHandler.ListJobApplications).Methods("GET")
	router.HandleFunc("/jobapplications", jobApplicationsHandler.CreateJobApplication).Methods("POST")
	router.HandleFunc("/jobapplications/{id}", jobApplicationsHandler.GetJobApplication).Methods("GET")
	router.HandleFunc("/jobapplications/{id}", jobApplicationsHandler.UpdateJobApplication).Methods("PUT")
	router.HandleFunc("/jobapplications/{id}", jobApplicationsHandler.DeleteJobApplication).Methods("DELETE")

	router.HandleFunc("/jobapplications", PreFlightHandler).Methods("OPTIONS")
	router.HandleFunc("/jobapplications/{id}", PreFlightHandler).Methods("OPTIONS")
	// Start server
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}


type jobApplicationStore interface {
	Add(id int, jobApplication jobApplications.JobApplication) error
	Get(id int) (*jobApplications.JobApplication, error)
	List() ([]jobApplications.JobApplication, error)
	Update(id int, jobApplication jobApplications.JobApplication) error
	Remove(id int) error
}

type JobApplicationsHandler struct {
	store jobApplicationStore
}

func NewJobApplicationHandler(s jobApplicationStore) *JobApplicationsHandler {
	return &JobApplicationsHandler{
		store: s,
	}
}

func (h JobApplicationsHandler) CreateJobApplication(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	var jobApplication jobApplications.JobApplication

	if err := json.NewDecoder(r.Body).Decode(&jobApplication); err != nil {
		print(fmt.Sprintf("Recieved following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	h.store.Add(jobApplication.Id, jobApplication)

	w.WriteHeader(http.StatusCreated)

}
func (h JobApplicationsHandler) ListJobApplications(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	jobapplications, err := jobApplicationStore.List(h.store)

	if err != nil {
		print(fmt.Sprintf("Recieved following error while getting jobApplications list: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}
	jsonBytes, err := json.Marshal(jobapplications)
	if len(jobapplications) == 0 {
		jsonBytes = []byte("[]")
	}

	if err != nil {
		InternalServerErrorHandler(w, r)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}
func (h JobApplicationsHandler) GetJobApplication(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	strId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("Recieved following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	jobApplication, err := h.store.Get(id)

	if err != nil {
		print(fmt.Sprintf("Recieved following error while getting jobApplication with id %d: \n%s", id, err.Error()))
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
	enableCors(&w)

	strId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("Recieved following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}
	oldJobApplication, err := jobApplicationStore.Get(h.store, id)
	if err != nil {
		print(fmt.Sprintf("Recieved following error while getting jobApplication with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	var newJobApplication jobApplications.JobApplication

	if err := json.NewDecoder(r.Body).Decode(&newJobApplication); err != nil {
		print(fmt.Sprintf("Recieved following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	newJobApplication.Id = oldJobApplication.Id

	err = h.store.Update(id, newJobApplication)

	if err != nil {
		print(fmt.Sprintf("Recieved following error while getting jobApplication with id %d: \n%s", id, err.Error()))
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
	enableCors(&w)

	strId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("Recieved following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	err = h.store.Remove(id)

	if err != nil {
		print(fmt.Sprintf("Recieved following error while getting jobApplication with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

func PreFlightHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.WriteHeader(http.StatusOK)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "kornelius.local")
	(*w).Header().Set("Access-Control-Request-Headers", "content-type")
	(*w).Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, OPTIONS")
	(*w).Header().Set("Accept", "/")
}
