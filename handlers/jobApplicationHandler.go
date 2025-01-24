package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/even44/JobsearchAPI/models"
	"github.com/even44/JobsearchAPI/jobapplicationstore"
	"github.com/gorilla/mux"
)

type JobApplicationsHandler struct {
	store  jobapplicationstore.JobApplicationStore
	logger *log.Logger
}

func NewJobApplicationHandler(s jobapplicationstore.JobApplicationStore) *JobApplicationsHandler {
	return &JobApplicationsHandler{
		store: s,
		logger: log.New(os.Stdout, "[JOBAPPLICATION HANDLER]", log.Ldate+log.Ltime+log.Lmsgprefix),
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

	resultJobApplication, err := h.store.AddApplication(jobApplication.Id, jobApplication)
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

	jobApplication, err := h.store.GetApplication(id)

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
	oldJobApplication, err := h.store.GetApplication(id)
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

	newJobApplication.Id = oldJobApplication.Id
	if newJobApplication.CompanyId == 0 {
		h.logger.Printf("[UPDATE] Updated job application had companyid = 0, using old company id: %d", oldJobApplication.CompanyId)
		newJobApplication.CompanyId = oldJobApplication.CompanyId
	}

	err = h.store.UpdateApplication(id, newJobApplication)

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

	err = h.store.RemoveApplication(id)

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

	var company models.Company
	h.logger.Printf("Received request to create company from: %s", r.Host)
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	resultCompany, err := h.store.AddCompany(company.Id, company)
	if err != nil {
		h.logger.Fatal(err.Error())
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

	h.logger.Printf("Received request to list companies from: %s", r.Host)
	companies, err := h.store.ListCompanies()
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
	h.logger.Printf("Received request to get company with id %s from: %s", strId, r.Host)
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
	h.logger.Printf("Received request to update company with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}
	oldCompany, err := h.store.GetCompany(id)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting company with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	var newCompany models.Company

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
	h.logger.Printf("Received request to delete company with id %s from: %s", strId, r.Host)
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

	var contact models.Contact
	h.logger.Printf("Received request to create contact from: %s", r.Host)
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	resultContact, err := h.store.AddContact(contact.CompanyId, contact)
	if err != nil {
		h.logger.Fatal(err.Error())
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

	h.logger.Printf("Received request to list contacts from: %s", r.Host)
	contacts, err := h.store.ListContacts()
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
	h.logger.Printf("Received request to get contact with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	contact, err := h.store.GetContact(id)

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
	h.logger.Printf("Received request to update contact with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}
	oldContact, err := h.store.GetContact(id)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting contact with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	var newContact models.Contact

	if err := json.NewDecoder(r.Body).Decode(&newContact); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err))
		InternalServerErrorHandler(w, r)
		return
	}

	h.logger.Printf("%d => %d", newContact.Id, oldContact.Id)
	newContact.Id = oldContact.Id

	err = h.store.UpdateContact(id, newContact)

	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting contact with id %d: \n%s", id, err))
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
	h.logger.Printf("Received request to delete contact with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	err = h.store.RemoveContact(id)

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
