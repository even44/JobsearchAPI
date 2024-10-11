package jobApplications

import "errors"

var (
	ErrNotFound = errors.New("not found")
)

type MemStore struct {
	list map[int]JobApplication
}

func NewMemStore() *MemStore {
	list := make(map[int]JobApplication)
	return &MemStore{
		list,
	}
}

func (m MemStore) Add(id int, recipe JobApplication) error {
	m.list[id] = recipe
	return nil
}

func (m MemStore) Get(id int) (JobApplication, error) {

	if val, ok := m.list[id]; ok {
		return val, nil
	}

	return JobApplication{}, ErrNotFound
}

func (m MemStore) List() ([]JobApplication, error) {
	var applicationList []JobApplication
	for jobApplication := range m.list {
		applicationList = append(applicationList, m.list[jobApplication])
	}

	return applicationList, nil
}

func (m MemStore) Update(id int, recipe JobApplication) error {

	if _, ok := m.list[id]; ok {
		m.list[id] = recipe
		return nil
	}

	return ErrNotFound
}

func (m MemStore) Remove(id int) error {
	delete(m.list, id)
	return nil
}