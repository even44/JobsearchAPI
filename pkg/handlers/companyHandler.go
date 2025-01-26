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

var CompanyH *CompanyHandler

type CompanyHandler struct {
	store  stores.CompanyStore
	logger *log.Logger
}

func NewCompanyHandler(s stores.CompanyStore) *CompanyHandler {
	return &CompanyHandler{
		store:  s,
		logger: log.New(os.Stdout, "[COMPANY] ", log.Ldate+log.Ltime+log.Lmsgprefix),
	}
}

func (h CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)

	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	h.logger.Printf("User id: %d", user.ID)

	var company models.Company
	company.UserID = user.ID
	h.logger.Printf("Received request to create company from: %s", r.Host)
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	resultCompany, err := h.store.AddCompany(company)
	if err != nil {
		h.logger.Println(err.Error())
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
func (h CompanyHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	if !checkOrigin(&w, r) {
		return
	}
	enableCors(&w)

	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	h.logger.Printf("User id: %d", user.ID)

	h.logger.Printf("Received request to list companies from: %s", r.Host)
	companies, err := h.store.ListCompanies(user.ID)
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
func (h CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
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

	company, err := h.store.GetCompany(uint(id))

	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting jobApplication with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	h.logger.Printf("User id: %d", user.ID)

	if company.UserID != user.ID {
		BadRequestHandler(w, r)
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
func (h CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
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

	oldCompany, err := h.store.GetCompany(uint(id))
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting company with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	h.logger.Printf("User id: %d", user.ID)

	if oldCompany.UserID != user.ID {
		BadRequestHandler(w, r)
		return
	}

	var newCompany models.Company

	if err := json.NewDecoder(r.Body).Decode(&newCompany); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	newCompany.ID = oldCompany.ID

	err = h.store.UpdateCompany(uint(id), newCompany)

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
func (h CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
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

	company, err := h.store.GetCompany(uint(id))

	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting jobApplication with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	h.logger.Printf("User id: %d", user.ID)

	if company.UserID != user.ID {
		BadRequestHandler(w, r)
		return
	}

	err = h.store.RemoveCompany(uint(id))

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
