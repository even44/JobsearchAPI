package jobApplications

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db_url = ""
var db_user = ""
var db_password = ""
var db_port = 0

var logger *log.Logger

func ParseEnv() {
	var temp string

	// Should look like "192.168.0.20" not "sixthousandandone"
	temp = os.Getenv("DB_URL")
	if temp != "" {
		db_url = temp
	}

	temp = os.Getenv("DB_USER")
	if temp != "" {
		db_user = temp
	}

	temp = os.Getenv("DB_PASSWORD")
	if temp != "" {
		db_password = temp
	}

	temp = os.Getenv("DB_PORT")
	if temp != "" {
		var err error
		db_port, err = strconv.Atoi(temp)
		if err != nil {
			log.Fatal("Could not convert DB_PORT to int")
			panic(err)
		}
	}

}

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
	err = db.AutoMigrate(&JobApplication{}, &Company{}, &Contact{})
	if err != nil {
		logger.Fatal(err)
	}

	return &MariaDBStore{
		db,
	}
}

func connectToMariaDB() (*gorm.DB, error) {
	ParseEnv()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/jobsearchdb?charset=utf8mb4&parseTime=True&loc=Local", db_user, db_password, db_url, db_port)
	logger.Printf("Connection string: %s", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s MariaDBStore) Add(id int, application JobApplication) (*JobApplication, error) {

	logger.Println("[ADD] Adding Company")
	company, err := s.AddCompany(0, application.Company)
	if err != nil {
		return nil, err
	}

	logger.Println("[ADD] Adding Contacts")
	for _, contact := range application.Company.Contacts {
		s.AddCompanyContact(company.Id, contact)
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
	resApplication, err := s.Get(application.Id)
	if err != nil {
		return nil, err
	}

	logger.Printf("Created job application with position '%s' and company '%s' with id %d ", resApplication.Position, resApplication.Company.Name, resApplication.Id)
	return resApplication, nil

}
func (s MariaDBStore) Get(Id int) (*JobApplication, error) {
	var application JobApplication
	result := s.db.Preload("Company.Contacts").Preload("Company").First(&application, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	logger.Printf("[GET] Got job application with position '%s' and company '%s' with id %d ", application.Position, application.Company.Name, application.Id)
	return &application, nil
}
func (s MariaDBStore) List() ([]JobApplication, error) {
	var applications []JobApplication
	s.db.Preload("Company.Contacts").Preload("Company").Find(&applications)
	logger.Printf("[LIST] Got %d job applications", len(applications))
	return applications, nil
}
func (s MariaDBStore) Update(id int, application JobApplication) error {
	result := s.db.Save(&application)
	if result.Error != nil {
		return result.Error
	}
	logger.Printf("[UPDATE] Updated job application with position '%s' and company '%s' with id %d ", application.Position, application.Company.Name, application.Id)
	return nil
}
func (s MariaDBStore) Remove(id int) error {
	application, err := s.Get(id)
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

func (s MariaDBStore) AddCompany(id int, company Company) (*Company, error) {

	err := s.db.First(&Company{}, Company{Name: company.Name, Location: company.Location}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Printf("[ADD] Company with name %s and location %s already exists and will not be created", company.Name, company.Location)
		s.db.First(&company, &Company{Name: company.Name, Location: company.Location})
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
func (s MariaDBStore) GetCompany(Id int) (*Company, error) {
	var company Company
	result := s.db.Preload("Contacts").First(&company, Id)
	if result.Error != nil {
		logger.Printf("[ERROR][GET] Could not find company with id %d", Id)
		return nil, result.Error
	}
	logger.Printf("[GET] Got company with Name '%s' and id %d ", company.Name, company.Id)
	return &company, nil
}
func (s MariaDBStore) ListCompanies() ([]Company, error) {
	var companies []Company
	s.db.Preload("Contacts").Find(&companies)
	logger.Printf("[LIST] Got %d companies", len(companies))
	return companies, nil
}
func (s MariaDBStore) UpdateCompany(Id int, company Company) error {
	var existingCompany *Company
	err := s.db.First(&existingCompany, Company{Name: company.Name, Location: company.Location}).Error
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

func (s MariaDBStore) AddCompanyContact(company_id int, contact Contact) (*Contact, error) {

	err := s.db.First(&Contact{}, Contact{Name: contact.Name}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Printf("[WARN][ADD] Contact with name %s already exists and will not be created", contact.Name)
		s.db.First(&contact, &Contact{Name: contact.Name})
		return &contact, err
	}
	contact.CompanyId = company_id
	result := s.db.Create(&contact)
	if result.Error != nil {
		return nil, result.Error
	}

	resConstact, err := s.GetCompanyContact(contact.Id)
	if err != nil {
		return nil, result.Error
	}
	logger.Printf("[ADD] Created contact with name %s", contact.Name)
	return resConstact, nil
}
func (s MariaDBStore) GetCompanyContact(id int) (*Contact, error) {
	var contact Contact
	result := s.db.First(&contact, id)
	if result.Error != nil {
		return nil, result.Error
	}
	logger.Printf("[GET] Got contact with Name '%s' and id %d ", contact.Name, contact.Id)
	return &contact, nil
}
func (s MariaDBStore) ListCompanyContacts() ([]Contact, error) {
	var contacts []Contact
	s.db.Find(&contacts)
	logger.Printf("[LIST] Got %d companies", len(contacts))
	return contacts, nil
}
func (s MariaDBStore) UpdateCompanyContract(id int, contact Contact) error {

	var existingContact *Contact
	err := s.db.First(&existingContact, Contact{Name: contact.Name}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if existingContact.Id != id {
			logger.Printf("[WARN][UPDATE] A different contact with name %s already exists and will not be updated", contact.Name)
			return err
		}
	}
	contact.CompanyId = existingContact.CompanyId
	result := s.db.Save(&contact)
	if result.Error != nil {
		logger.Printf("[ERROR] %s", result.Error)
		return result.Error
	}
	logger.Printf("[UPDATE] Updated contact with name '%s' and company_id '%d' with id %d ", contact.Name, contact.CompanyId, contact.Id)
	return nil
}
func (s MariaDBStore) RemoveCompanyContact(id int) error {
	contact, err := s.GetCompanyContact(id)
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
