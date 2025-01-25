package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/even44/JobsearchAPI/pkg/handlers"
	"github.com/even44/JobsearchAPI/pkg/initializers"
	"github.com/even44/JobsearchAPI/pkg/middleware"
	"github.com/gorilla/mux"
)

var logger *log.Logger

func init() {

	logger = log.New(os.Stdout, "MAIN: ", log.Ldate+log.Ltime+log.Lmsgprefix)

	initializers.LoadEnvVariables()
	initializers.ParseEnvVariables()
	initializers.ConnectToMariaDB()
	initializers.SyncDatabase()

	handlers.UH = handlers.NewUserHandler(initializers.Store)
	handlers.JAH = handlers.NewJobApplicationHandler(initializers.Store)
}

func main() {

	// Create the store and Jobapplication handler

	// Create public
	global := mux.NewRouter().NewRoute()
	auth := global.PathPrefix("/auth").Subrouter()
	public := global.PathPrefix("/public").Subrouter()

	auth.HandleFunc("/jobapplications", handlers.JAH.ListJobApplications).Methods("GET")
	auth.HandleFunc("/jobapplications", handlers.JAH.CreateJobApplication).Methods("POST")
	auth.HandleFunc("/jobapplications/{id}", handlers.JAH.GetJobApplication).Methods("GET")
	auth.HandleFunc("/jobapplications/{id}", handlers.JAH.UpdateJobApplication).Methods("PUT")
	auth.HandleFunc("/jobapplications/{id}", handlers.JAH.DeleteJobApplication).Methods("DELETE")
	auth.HandleFunc("/jobapplications", handlers.PreFlightHandler).Methods("OPTIONS")
	auth.HandleFunc("/jobapplications/{id}", handlers.PreFlightHandler).Methods("OPTIONS")

	auth.HandleFunc("/companies", handlers.JAH.ListCompanies).Methods("GET")
	auth.HandleFunc("/companies", handlers.JAH.CreateCompany).Methods("POST")
	auth.HandleFunc("/companies/{id}", handlers.JAH.GetCompany).Methods("GET")
	auth.HandleFunc("/companies/{id}", handlers.JAH.UpdateCompany).Methods("PUT")
	auth.HandleFunc("/companies/{id}", handlers.JAH.DeleteCompany).Methods("DELETE")
	auth.HandleFunc("/companies", handlers.PreFlightHandler).Methods("OPTIONS")
	auth.HandleFunc("/companies/{id}", handlers.PreFlightHandler).Methods("OPTIONS")

	auth.HandleFunc("/contacts", handlers.JAH.ListContacts).Methods("GET")
	auth.HandleFunc("/contacts", handlers.JAH.CreateContact).Methods("POST")
	auth.HandleFunc("/contacts/{id}", handlers.JAH.GetContact).Methods("GET")
	auth.HandleFunc("/contacts/{id}", handlers.JAH.UpdateContact).Methods("PUT")
	auth.HandleFunc("/contacts/{id}", handlers.JAH.DeleteContact).Methods("DELETE")
	auth.HandleFunc("/contacts", handlers.PreFlightHandler).Methods("OPTIONS")
	auth.HandleFunc("/contacts/{id}", handlers.PreFlightHandler).Methods("OPTIONS")

	public.HandleFunc("/signup", handlers.UH.SignUp).Methods("POST")
	public.HandleFunc("/login", handlers.UH.Login).Methods("POST")

	public.Use(middleware.Logging)

	// Start server
	logger.Printf("Jobsearch API running on port: %d\n", initializers.ApiPort)
	http.ListenAndServe(fmt.Sprintf(":%d", initializers.ApiPort), nil)
}
