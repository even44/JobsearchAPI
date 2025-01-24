package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/even44/JobsearchAPI/pkg/jobApplications"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var port int = 3001
var trusted_origin = ""
var logger *log.Logger

func main() {

	logger = log.New(os.Stdout, "MAIN: ", log.Ldate+log.Ltime+log.Lmsgprefix)

	logger.Println("Loading .env file")
	err := godotenv.Load()
	if err != nil {
		logger.Println("[ERROR] No .env file or error loading, skipping")
	}

	ParseEnv()

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

	router.HandleFunc("/companies", jobApplicationsHandler.ListCompanies).Methods("GET")
	router.HandleFunc("/companies", jobApplicationsHandler.CreateCompany).Methods("POST")
	router.HandleFunc("/companies/{id}", jobApplicationsHandler.GetCompany).Methods("GET")
	router.HandleFunc("/companies/{id}", jobApplicationsHandler.UpdateCompany).Methods("PUT")
	router.HandleFunc("/companies/{id}", jobApplicationsHandler.DeleteCompany).Methods("DELETE")

	router.HandleFunc("/contacts", jobApplicationsHandler.ListContacts).Methods("GET")
	router.HandleFunc("/contacts", jobApplicationsHandler.CreateContact).Methods("POST")
	router.HandleFunc("/contacts/{id}", jobApplicationsHandler.GetContact).Methods("GET")
	router.HandleFunc("/contacts/{id}", jobApplicationsHandler.UpdateContact).Methods("PUT")
	router.HandleFunc("/contacts/{id}", jobApplicationsHandler.DeleteContact).Methods("DELETE")

	router.HandleFunc("/jobapplications", PreFlightHandler).Methods("OPTIONS")
	router.HandleFunc("/jobapplications/{id}", PreFlightHandler).Methods("OPTIONS")
	router.HandleFunc("/companies", PreFlightHandler).Methods("OPTIONS")
	router.HandleFunc("/companies/{id}", PreFlightHandler).Methods("OPTIONS")
	router.HandleFunc("/contacts", PreFlightHandler).Methods("OPTIONS")
	router.HandleFunc("/contacts/{id}", PreFlightHandler).Methods("OPTIONS")
	// Start server
	logger.Printf("Jobsearch API running on port: %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

func ParseEnv() {
	logger.Println("Getting API env variables")
	var temp string

	// Should look like "6001" not "sixthousandandone"
	temp = os.Getenv("API_PORT")
	if temp != "" {
		var err error
		port, err = strconv.Atoi(temp)
		if err != nil {
			fmt.Println("[ERROR] Could not convert API_PORT to int")
			panic(err)
		}
	}

	// Should look like "http://ip:port" or "https://domain.example"
	temp = os.Getenv("TRUSTED_ORIGIN")
	if temp != "" {
		trusted_origin = temp
	}
}

type jobApplicationStore interface {
	Add(id int, jobApplication jobApplications.JobApplication) (*jobApplications.JobApplication, error)
	Get(id int) (*jobApplications.JobApplication, error)
	List() ([]jobApplications.JobApplication, error)
	Update(id int, jobApplication jobApplications.JobApplication) error
	Remove(id int) error
	AddCompany(id int, company jobApplications.Company) (*jobApplications.Company, error)
	GetCompany(Id int) (*jobApplications.Company, error)
	ListCompanies() ([]jobApplications.Company, error)
	UpdateCompany(Id int, company jobApplications.Company) error
	RemoveCompany(id int) error
	AddCompanyContact(company_id int, contact jobApplications.Contact) (*jobApplications.Contact, error)
	GetCompanyContact(id int) (*jobApplications.Contact, error)
	ListCompanyContacts() ([]jobApplications.Contact, error)
	UpdateCompanyContract(id int, contact jobApplications.Contact) error
	RemoveCompanyContact(id int) error
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
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)

	var jobApplication jobApplications.JobApplication
	logger.Printf("Received request to create job application from: %s", r.Host)
	if err := json.NewDecoder(r.Body).Decode(&jobApplication); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	resultJobApplication, err := h.store.Add(jobApplication.Id, jobApplication)
	if err != nil {
		logger.Println(err.Error())
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

	logger.Printf("Received request to list job applications from: %s", r.Host)
	jobapplications, err := jobApplicationStore.List(h.store)
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
	logger.Printf("Received request to get job application with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	jobApplication, err := h.store.Get(id)

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
	logger.Printf("Received request to update job application with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}
	oldJobApplication, err := jobApplicationStore.Get(h.store, id)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting jobApplication with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	var newJobApplication jobApplications.JobApplication

	if err := json.NewDecoder(r.Body).Decode(&newJobApplication); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	newJobApplication.Id = oldJobApplication.Id
	if newJobApplication.CompanyId == 0 {
		logger.Printf("[UPDATE] Updated job application had companyid = 0, using old company id: %d", oldJobApplication.CompanyId)
		newJobApplication.CompanyId = oldJobApplication.CompanyId
	}

	err = h.store.Update(id, newJobApplication)

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
	logger.Printf("Received request to delete job application with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	err = h.store.Remove(id)

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

func (h JobApplicationsHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)

	var company jobApplications.Company
	logger.Printf("Received request to create company from: %s", r.Host)
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	resultCompany, err := h.store.AddCompany(company.Id, company)
	if err != nil {
		logger.Fatal(err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(resultCompany)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)

}
func (h JobApplicationsHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)

	logger.Printf("Received request to list companies from: %s", r.Host)
	companies, err := jobApplicationStore.ListCompanies(h.store)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting company list: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}
	jsonBytes, err := json.Marshal(companies)
	if len(companies) == 0 {
		jsonBytes = []byte("[]")
	}

	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}
func (h JobApplicationsHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)
	strId := mux.Vars(r)["id"]
	logger.Printf("Received request to get company with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	company, err := h.store.GetCompany(id)

	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting jobApplication with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(company)

	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h JobApplicationsHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)

	strId := mux.Vars(r)["id"]
	logger.Printf("Received request to update company with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}
	oldCompany, err := jobApplicationStore.GetCompany(h.store, id)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting company with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	var newCompany jobApplications.Company

	if err := json.NewDecoder(r.Body).Decode(&newCompany); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	newCompany.Id = oldCompany.Id

	err = h.store.UpdateCompany(id, newCompany)

	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while updating jobApplication with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h JobApplicationsHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)
	strId := mux.Vars(r)["id"]
	logger.Printf("Received request to delete company with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	err = h.store.RemoveCompany(id)

	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting company with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h JobApplicationsHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)

	var contact jobApplications.Contact
	logger.Printf("Received request to create contact from: %s", r.Host)
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	resultContact, err := h.store.AddCompanyContact(contact.CompanyId, contact)
	if err != nil {
		logger.Fatal(err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(resultContact)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)

}
func (h JobApplicationsHandler) ListContacts(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)

	logger.Printf("Received request to list contacts from: %s", r.Host)
	contacts, err := jobApplicationStore.ListCompanyContacts(h.store)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting contact list: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}
	jsonBytes, err := json.Marshal(contacts)
	if len(contacts) == 0 {
		jsonBytes = []byte("[]")
	}

	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}
func (h JobApplicationsHandler) GetContact(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)
	strId := mux.Vars(r)["id"]
	logger.Printf("Received request to get contact with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	contact, err := h.store.GetCompanyContact(id)

	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting contact with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(contact)

	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h JobApplicationsHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)

	strId := mux.Vars(r)["id"]
	logger.Printf("Received request to update contact with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}
	oldContact, err := jobApplicationStore.GetCompanyContact(h.store, id)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting contact with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	var newContact jobApplications.Contact

	if err := json.NewDecoder(r.Body).Decode(&newContact); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	logger.Printf("%d => %d", newContact.Id, oldContact.Id)
	newContact.Id = oldContact.Id

	err = h.store.UpdateCompanyContract(id, newContact)

	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting contact with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h JobApplicationsHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)
	strId := mux.Vars(r)["id"]
	logger.Printf("Received request to delete contact with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	err = h.store.RemoveCompanyContact(id)

	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting contact with id %d: \n%s", id, err.Error()))
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
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)
	w.WriteHeader(http.StatusOK)
}
func checkOrigin(w *http.ResponseWriter, r *http.Request) bool {
	if r.Header.Get("Origin") == trusted_origin {
		(*w).Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		return true
	} else {
		logger.Printf("Received request with wrong origin, %s, from: %s", r.Header.Get("Origin"), r.Host)
		InternalServerErrorHandler((*w), r)
		return false
	}
}
func enableCors(w *http.ResponseWriter) {

	(*w).Header().Set("Access-Control-Allow-Headers", "content-type")
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, OPTIONS, DELETE")

}
