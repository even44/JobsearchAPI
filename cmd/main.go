package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/even44/JobsearchAPI/pkg/jobApplications"
	"github.com/gorilla/mux"
)

func main() {

	const port int = 3001

	fmt.Printf("Jobsearch API running on port: %d\n", port)

	// Create the store and Jobapplication handler
	store := jobApplications.NewMemStore()
	jobApplicationsHandler := NewJobApplicationHandler(store)
	home := homeHandler{}

	// Create router
	router := mux.NewRouter()

	router.HandleFunc("/", home.ServeHTTP)
	router.HandleFunc("/jobapplications", jobApplicationsHandler.ListJobApplications).Methods("GET")
	router.HandleFunc("/jobapplications", jobApplicationsHandler.CreateJobApplication).Methods("POST")
	router.HandleFunc("/jobapplications{id}", jobApplicationsHandler.GetJobApplication).Methods("GET")
	router.HandleFunc("/jobapplications{id}", jobApplicationsHandler.UpdateJobApplication).Methods("PUT")
	router.HandleFunc("/jobapplications{id}", jobApplicationsHandler.DeleteJobApplication).Methods("DELETE")

	// Start server
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page i guess"))
}

type jobApplicationStore interface {
	Add(id int, jobApplication jobApplications.JobApplication) error
	Get(id int) (jobApplications.JobApplication, error)
	List() (map[int]jobApplications.JobApplication, error)
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
	var jobApplication jobApplications.JobApplication
	
	if err := json.NewDecoder(r.Body).Decode(&jobApplication); err != nil {
		print(fmt.Sprintf("Recieved following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}
	
	h.store.Add(jobApplication.Id, jobApplication)
	
	w.WriteHeader(http.StatusCreated)
	
}
func (h JobApplicationsHandler) ListJobApplications(w http.ResponseWriter, r *http.Request)  {
	
	jobapplications, err := jobApplicationStore.List(h.store)
	
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	jsonBytes, err := json.Marshal(jobapplications)
	
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h JobApplicationsHandler) GetJobApplication(w http.ResponseWriter, r *http.Request)    {}
func (h JobApplicationsHandler) UpdateJobApplication(w http.ResponseWriter, r *http.Request) {}
func (h JobApplicationsHandler) DeleteJobApplication(w http.ResponseWriter, r *http.Request) {}


func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}