package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/even44/JobsearchAPI/handlers"
	"github.com/even44/JobsearchAPI/initializers"
	"github.com/even44/JobsearchAPI/jobapplicationstore"
	"github.com/gorilla/mux"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "MAIN: ", log.Ldate+log.Ltime+log.Lmsgprefix)
	initializers.LoadEnvVariables()
	initializers.ParseEnvVariables()
}

func main() {

	// Create the store and Jobapplication handler
	store := jobapplicationstore.NewMariaDBStore()
	jobApplicationHandler := handlers.NewJobApplicationHandler(store)

	// Create router
	router := mux.NewRouter()

	router.HandleFunc("/jobapplications", jobApplicationHandler.ListJobApplications).Methods("GET")
	router.HandleFunc("/jobapplications", jobApplicationHandler.CreateJobApplication).Methods("POST")
	router.HandleFunc("/jobapplications/{id}", jobApplicationHandler.GetJobApplication).Methods("GET")
	router.HandleFunc("/jobapplications/{id}", jobApplicationHandler.UpdateJobApplication).Methods("PUT")
	router.HandleFunc("/jobapplications/{id}", jobApplicationHandler.DeleteJobApplication).Methods("DELETE")

	router.HandleFunc("/companies", jobApplicationHandler.ListCompanies).Methods("GET")
	router.HandleFunc("/companies", jobApplicationHandler.CreateCompany).Methods("POST")
	router.HandleFunc("/companies/{id}", jobApplicationHandler.GetCompany).Methods("GET")
	router.HandleFunc("/companies/{id}", jobApplicationHandler.UpdateCompany).Methods("PUT")
	router.HandleFunc("/companies/{id}", jobApplicationHandler.DeleteCompany).Methods("DELETE")

	router.HandleFunc("/contacts", jobApplicationHandler.ListContacts).Methods("GET")
	router.HandleFunc("/contacts", jobApplicationHandler.CreateContact).Methods("POST")
	router.HandleFunc("/contacts/{id}", jobApplicationHandler.GetContact).Methods("GET")
	router.HandleFunc("/contacts/{id}", jobApplicationHandler.UpdateContact).Methods("PUT")
	router.HandleFunc("/contacts/{id}", jobApplicationHandler.DeleteContact).Methods("DELETE")

	router.HandleFunc("/jobapplications", handlers.PreFlightHandler).Methods("OPTIONS")
	router.HandleFunc("/jobapplications/{id}", handlers.PreFlightHandler).Methods("OPTIONS")
	router.HandleFunc("/companies", handlers.PreFlightHandler).Methods("OPTIONS")
	router.HandleFunc("/companies/{id}", handlers.PreFlightHandler).Methods("OPTIONS")
	router.HandleFunc("/contacts", handlers.PreFlightHandler).Methods("OPTIONS")
	router.HandleFunc("/contacts/{id}", handlers.PreFlightHandler).Methods("OPTIONS")
	// Start server
	logger.Printf("Jobsearch API running on port: %d\n", initializers.ApiPort)
	http.ListenAndServe(fmt.Sprintf(":%d", initializers.ApiPort), router)
}
