package handlers

import (
	"github.com/goccy/go-json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/even44/JobsearchAPI/pkg/models"
	"github.com/even44/JobsearchAPI/pkg/stores"
	"github.com/gorilla/mux"
)

var ContactH *ContactHandler

type ContactHandler struct {
	store  stores.ContactStore
	logger *log.Logger
}

func NewContactHandler(s stores.ContactStore) *ContactHandler {
	return &ContactHandler{
		store:  s,
		logger: log.New(os.Stdout, "[CONTACT] ", log.Ldate+log.Ltime+log.Lmsgprefix),
	}
}

func (h ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {

	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	h.logger.Printf("User id: %d", user.ID)

	var contact models.Contact
	contact.UserID = user.ID
	h.logger.Printf("Received request to create contact from: %s", r.Host)
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	resultContact, err := h.store.AddContact(contact)
	if err != nil {
		h.logger.Println(err.Error())
		BadRequestHandler(w, r)
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
func (h ContactHandler) ListContacts(w http.ResponseWriter, r *http.Request) {


	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	h.logger.Printf("User id: %d", user.ID)

	h.logger.Printf("Received request to list contacts from: %s", r.Host)
	contacts, err := h.store.ListContacts(user.ID)
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
func (h ContactHandler) GetContact(w http.ResponseWriter, r *http.Request) {

	strId := mux.Vars(r)["id"]
	h.logger.Printf("Received request to get contact with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	contact, err := h.store.GetContact(uint(id))

	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting contact with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	h.logger.Printf("User id: %d", user.ID)

	if contact.UserID != user.ID {
		BadRequestHandler(w, r)
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
func (h ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {


	strId := mux.Vars(r)["id"]
	h.logger.Printf("Received request to update contact with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}
	oldContact, err := h.store.GetContact(uint(id))
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting contact with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	h.logger.Printf("User id: %d", user.ID)

	if oldContact.UserID != user.ID {
		BadRequestHandler(w, r)
		return
	}

	var newContact models.Contact

	if err := json.NewDecoder(r.Body).Decode(&newContact); err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while parsing request JSON: \n%s", err))
		InternalServerErrorHandler(w, r)
		return
	}

	h.logger.Printf("%d => %d", newContact.ID, oldContact.ID)
	newContact.ID = oldContact.ID

	err = h.store.UpdateContact(uint(id), newContact)

	if err != nil {
		h.logger.Printf(
			"[ERROR] Received following error while getting contact with id %d: \n%s", id, err)
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {

	strId := mux.Vars(r)["id"]
	h.logger.Printf("Received request to delete contact with id %s from: %s", strId, r.Host)
	id, err := strconv.Atoi(strId)
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while converting id to int \n%s", err.Error()))
		InternalServerErrorHandler(w, r)
		return
	}

	var user *models.User = r.Context().Value(models.User{}).(*models.User)
	h.logger.Printf("User id: %d", user.ID)

	contact, err := h.store.GetContact(uint(id))
	if err != nil {
		print(fmt.Sprintf("[ERROR] Received following error while getting contact with id %d: \n%s", id, err.Error()))
		if err.Error() == "not found" {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}

	if contact.UserID != user.ID {
		BadRequestHandler(w, r)
		return
	}


	err = h.store.RemoveContact(uint(id))

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
