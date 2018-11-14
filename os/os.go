package os

import (
	"os"

	"github.com/go-franky/keys"
)

type oskm struct{}

// NewFromOS returns a manager that uses environment
// variables to read / write from to
func NewFromOS() keys.Manager {
	return &oskm{}
}

func (km *oskm) Lookup(k string) (string, bool) {
	return os.LookupEnv(k)
}

func (km *oskm) Set(k, v string) error {
	return os.Setenv(k, v)
}

func (km *oskm) Get(key string) string {
	k, _ := km.Lookup(key)
	return k
}
