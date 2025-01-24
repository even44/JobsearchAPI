package jobapplicationstore

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/even44/JobsearchAPI/initializers"
	"github.com/even44/JobsearchAPI/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var logger *log.Logger

type MariaDBStore struct {
	db *gorm.DB
}

func NewMariaDBStore() *MariaDBStore {

	logger = log.New(os.Stdout, "MariaDBStore: ", log.Ldate+log.Ltime+log.Lmsgprefix)

	db, err := connectToMariaDB()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("Auto migrating database")
	err = db.AutoMigrate(&models.JobApplication{}, &models.Company{}, &models.Contact{})
	if err != nil {
		logger.Fatal(err)
	}

	return &MariaDBStore{
		db,
	}
}

func connectToMariaDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/jobsearchdb?charset=utf8mb4&parseTime=True&loc=Local",
		initializers.DbUser, initializers.DbPassword, initializers.DbURL, initializers.DbPort)
	logger.Printf("Connection string: %s", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s MariaDBStore) AddApplication(id int, application models.JobApplication) (*models.JobApplication, error) {

	logger.Println("[ADD] Adding Company")
	company, err := s.AddCompany(0, application.Company)
	if err != nil {
		return nil, err
	}

	logger.Println("[ADD] Adding Contacts")
	for _, contact := range application.Company.Contacts {
		s.AddContact(company.Id, contact)
	}
	logger.Println("[ADD] Done Adding Contacts")

	logger.Println("[ADD] Creating Application")
	logger.Printf("[ADD] Application id %d => %d", application.CompanyId, company.Id)
	application.CompanyId = company.Id
	result := s.db.Omit("Company").Create(&application)
	if result.Error != nil {
		return nil, result.Error
	}

	logger.Println("[ADD] Getting created application")
	resApplication, err := s.GetApplication(application.Id)
	if err != nil {
		return nil, err
	}

	logger.Printf("[ADD] Created job application with position '%s' and company '%s' with id %d ", resApplication.Position, resApplication.Company.Name, resApplication.Id)
	return resApplication, nil

}
func (s MariaDBStore) GetApplication(Id int) (*models.JobApplication, error) {
	var application models.JobApplication
	result := s.db.Preload("Company.Contacts").Preload("Company").First(&application, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	logger.Printf("[GET] Got job application with position '%s' and company '%s' with id %d ", application.Position, application.Company.Name, application.Id)
	return &application, nil
}
func (s MariaDBStore) ListApplications() ([]models.JobApplication, error) {
	var applications []models.JobApplication
	s.db.Preload("Company.Contacts").Preload("Company").Find(&applications)
	logger.Printf("[LIST] Got %d job applications", len(applications))
	return applications, nil
}
func (s MariaDBStore) UpdateApplication(id int, application models.JobApplication) error {
	result := s.db.Omit("Company").Save(&application)
	if result.Error != nil {
		return result.Error
	}
	logger.Printf("[UPDATE] Updated job application with position '%s' with id %d ", application.Position, application.Id)
	return nil
}
func (s MariaDBStore) RemoveApplication(id int) error {
	application, err := s.GetApplication(id)
	if err != nil {
		return err
	}
	result := s.db.Delete(application)
	if result.Error != nil {
		return result.Error
	}
	logger.Printf("[DELETE] Deleted job application with position '%s' and company '%s' with id %d ", application.Position, application.Company.Name, application.Id)
	return nil
}

func (s MariaDBStore) AddCompany(id int, company models.Company) (*models.Company, error) {

	err := s.db.First(&models.Company{}, models.Company{Name: company.Name, Location: company.Location}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Printf("[WARN][ADD] Company with name %s and location %s already exists and will not be created", company.Name, company.Location)
		s.db.First(&company, &models.Company{Name: company.Name, Location: company.Location})
		return &company, err
	}

	logger.Printf("[ADD] Company with name %s and location %s does not exist and will be created", company.Name, company.Location)
	result := s.db.Omit("Contacts").Create(&company)
	if result.Error != nil {
		return &company, result.Error
	}

	resCompany, err := s.GetCompany(company.Id)
	if err != nil {
		return nil, result.Error
	}
	logger.Printf("[ADD] Created company with Name '%s' and id %d ", company.Name, company.Id)
	return resCompany, nil
}
func (s MariaDBStore) GetCompany(Id int) (*models.Company, error) {
	var company models.Company
	result := s.db.Preload("Contacts").First(&company, Id)
	if result.Error != nil {
		logger.Printf("[ERROR][GET] Could not find company with id %d", Id)
		return nil, result.Error
	}
	logger.Printf("[GET] Got company with Name '%s' and id %d ", company.Name, company.Id)
	return &company, nil
}
func (s MariaDBStore) ListCompanies() ([]models.Company, error) {
	var companies []models.Company
	s.db.Preload("Contacts").Find(&companies)
	logger.Printf("[LIST] Got %d companies", len(companies))
	return companies, nil
}
func (s MariaDBStore) UpdateCompany(Id int, company models.Company) error {
	var existingCompany *models.Company
	err := s.db.First(&existingCompany, models.Company{Name: company.Name, Location: company.Location}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if existingCompany.Id != Id {
			logger.Printf("[WARN][UPDATE] Company with name %s and location %s already exists and will not be updated", company.Name, company.Location)
			return err
		}
	}

	result := s.db.Save(&company)
	if result.Error != nil {
		logger.Printf("[ERROR] %s", result.Error)
		return result.Error
	}
	logger.Printf("[UPDATE] Updated company with name '%s' and location '%s' with id %d ", company.Name, company.Location, company.Id)
	return nil
}
func (s MariaDBStore) RemoveCompany(id int) error {
	company, err := s.GetCompany(id)
	if err != nil {
		return err
	}
	result := s.db.Delete(company)
	if result.Error != nil {
		return result.Error
	}
	logger.Printf("[DELETE] Deleted company with name '%s' and location '%s' with id %d ", company.Name, company.Location, company.Id)
	return nil
}

func (s MariaDBStore) AddContact(company_id int, contact models.Contact) (*models.Contact, error) {

	err := s.db.First(&models.Contact{}, models.Contact{Name: contact.Name}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Printf("[WARN][ADD] Contact with name %s already exists and will not be created", contact.Name)
		s.db.First(&contact, &models.Contact{Name: contact.Name})
		return &contact, err
	}
	contact.CompanyId = company_id
	result := s.db.Create(&contact)
	if result.Error != nil {
		return nil, result.Error
	}

	resConstact, err := s.GetContact(contact.Id)
	if err != nil {
		return nil, result.Error
	}
	logger.Printf("[ADD] Created contact with name %s", contact.Name)
	return resConstact, nil
}
func (s MariaDBStore) GetContact(id int) (*models.Contact, error) {
	var contact models.Contact
	result := s.db.First(&contact, id)
	if result.Error != nil {
		return nil, result.Error
	}
	logger.Printf("[GET] Got contact with Name '%s' and id %d ", contact.Name, contact.Id)
	return &contact, nil
}
func (s MariaDBStore) ListContacts() ([]models.Contact, error) {
	var contacts []models.Contact
	s.db.Find(&contacts)
	logger.Printf("[LIST] Got %d companies", len(contacts))
	return contacts, nil
}
func (s MariaDBStore) UpdateContact(id int, contact models.Contact) error {

	var existingContact *models.Contact
	err := s.db.First(&existingContact, models.Contact{Name: contact.Name}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if existingContact.Id != id {
			logger.Printf("[WARN][UPDATE] A different contact with name %s already exists and will not be updated", contact.Name)
			return err
		}
	}

	if contact.CompanyId == 0 {
		contact.CompanyId = existingContact.CompanyId
	}

	result := s.db.Save(&contact)
	if result.Error != nil {
		logger.Printf("[ERROR] %s", result.Error)
		return result.Error
	}
	logger.Printf("[UPDATE] Updated contact with name '%s' and company_id '%d' with id %d ", contact.Name, contact.CompanyId, contact.Id)
	return nil
}
func (s MariaDBStore) RemoveContact(id int) error {
	contact, err := s.GetContact(id)
	if err != nil {
		return err
	}
	result := s.db.Delete(contact)
	if result.Error != nil {
		return result.Error
	}
	logger.Printf("[DELETE] Deleted contact with name '%s' and company_id '%d' with id %d ", contact.Name, contact.CompanyId, contact.Id)
	return nil
}
