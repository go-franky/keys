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

func TestMultiManager(t *testing.T) {
	defaultMgr := keys.NewKeyManager()
	empty := &emptyMgr{}
	all := []keys.Manager{defaultMgr, empty}

	var tAll = keys.MultiManager(all...)

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
	tAll = keys.MultiManager(defaultMgr, &T1{})
	if tAll.Get("Hello") != "World" {
		t.Errorf("expected %v, got: %v", "World", tAll.Get("Hello"))
	}

	// Check the order
	defaultMgr.Set("Hello", "This is my World")
	if tAll.Get("Hello") != "This is my World" {
		t.Errorf("expected %v, got: %v", "This is My World", tAll.Get("Hello"))
	}
}

func TestMultiGetter(t *testing.T) {
	defaultMgr := keys.NewKeyManager()
	empty := &emptyMgr{}
	all := []keys.Getter{defaultMgr, empty}

	var tAll = keys.MultiGetter(all...)

	if k := tAll.Get("Non-Exisitent"); k != "" {
		t.Errorf("expected empty, got %v", k)
	}

	all = []keys.Getter{empty, &T1{}}
	tAll = keys.MultiGetter(all...)

	if k := tAll.Get("Hello"); k != "World" {
		t.Errorf("expected %v, got %v", "World", k)
	}
}

func TestMultiLookuper(t *testing.T) {
	defaultMgr := keys.NewKeyManager()
	empty := &emptyMgr{}
	all := []keys.Lookuper{defaultMgr, empty}

	var tAll = keys.MultiLookuper(all...)

	if k, ok := tAll.Lookup("Non-Existient"); k != "" || ok {
		t.Errorf("expected a combined manager to be initialized empty for lookups, got %v - %v", k, ok)
	}

	all = []keys.Lookuper{empty, &T1{}}
	tAll = keys.MultiLookuper(all...)

	if k, ok := tAll.Lookup("Hello"); k != "World" || !ok {
		t.Errorf("expected a multi lookup to retur right value for key: %v, bool: %v; got key: %v, bool: %v", "World", true, k, ok)
	}
}
