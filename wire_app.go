package learning

import (
	"strconv"
)

type repo interface {
	Get(id int) string
}

type memRepo struct{}

func newMemRepo() *memRepo {
	return &memRepo{}
}

func (m *memRepo) Get(id int) string {
	return strconv.Itoa(id)
}

type jsonRepo struct{}

func newJsonRepo() *jsonRepo {
	return &jsonRepo{}
}

func (m *jsonRepo) Get(id int) string {
	return `{"id":` + strconv.Itoa(id) + `}`
}

type app struct {
	memRepo  repo
	jsonRepo repo
}

func newApp(memRepo repo, jsonRepo repo) *app {
	return &app{
		memRepo:  memRepo,
		jsonRepo: jsonRepo,
	}
}
