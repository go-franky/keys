package keys_test

import (
	"testing"

	"github.com/go-franky/keys"
)

type T1 struct{}

func (km *T1) Get(k string) string            { return "World" }
func (km *T1) Lookup(k string) (string, bool) { return "World", true }
func (km *T1) Set(k, v string) error          { return nil }

type emptyMgr struct{}

func (km *emptyMgr) Get(k string) string            { return "" }
func (km *emptyMgr) Lookup(k string) (string, bool) { return "", false }
func (km *emptyMgr) Set(k, v string) error          { return nil }

func TestManager(t *testing.T) {
	km := keys.NewKeyManager()
	if v := km.Get("Hello"); v != "" {
		t.Errorf("expected a manager to be initialized empty got %v", v)
	}

	if k, ok := km.Lookup("Hello"); k != "" || ok {
		t.Errorf("expected a manager to be initialized empty for lookups, got %v - %v", k, ok)
	}

	km.Set("Hello", "World")

	if k := km.Get("Hello"); k != "World" {
		t.Errorf("expected %v, got %v", "World", k)
	}

	if k, ok := km.Lookup("Hello"); k != "World" || !ok {
		t.Errorf("expected a manager to be initialized empty for lookups, got %v - %v", k, ok)
	}
}

func TestCombine(t *testing.T) {
	defaultMgr := keys.NewKeyManager()
	empty := &emptyMgr{}
	all := []keys.Manager{defaultMgr, empty}

	var tAll = keys.Combine(all...)

	if k, ok := tAll.Lookup("Non-Existient"); k != "" || ok {
		t.Errorf("expected a combined manager to be initialized empty for lookups, got %v - %v", k, ok)
	}
	if k := tAll.Get("Non-Exisitent"); k != "" {
		t.Errorf("expected empty, got %v", k)
	}

	// Setting tAll does not change the others
	tAll.Set("Hey", "You")

	if k, ok := tAll.Lookup("Hey"); k != "You" || !ok {
		t.Errorf("expected a combined manager to be initialized empty for lookups, got %v - %v", "You", ok)
	}

	if k := tAll.Get("Hey"); k != "You" {
		t.Errorf("expected %v, got %v", "You", k)
	}

	for _, mgr := range all {
		if k, ok := mgr.Lookup("Non-Existient"); k != "" || ok {
			t.Errorf("expected a combined manager to be initialized empty for lookups, got %v - %v", k, ok)
		}
		if k := mgr.Get("Non-Exisitent"); k != "" {
			t.Errorf("expected empty, got %v", k)
		}
	}

	// Reading any of the listed ones
	tAll = keys.Combine(defaultMgr, &T1{})
	if tAll.Get("Hello") != "World" {
		t.Errorf("expected %v, got: %v", "World", tAll.Get("Hello"))
	}

	// Check the order
	defaultMgr.Set("Hello", "This is my World")
	if tAll.Get("Hello") != "This is my World" {
		t.Errorf("expected %v, got: %v", "This is My World", tAll.Get("Hello"))
	}
}
