package python

import (
	"log"

	py "github.com/sbinet/go-python"
)

type python struct {
	modules map[string]*py.PyObject
}

func New() {
	if er := python.Initilaize(); err != nil {
		log.Fatal("dead")
	}
}
