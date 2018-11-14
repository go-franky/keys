package os_test

import (
	"os"
	"testing"

	osMgr "github.com/go-franky/keys/os"
)

func TestOSMgr(t *testing.T) {
	osMgr := osMgr.NewFromOS()

	if osMgr.Get("HOME") != os.Getenv("HOME") {
		t.Errorf("expected mgr: %v to eql env %v", osMgr.Get("HOME"), os.Getenv("HOME"))
	}
	osMgr.Set("Hello", "World")
	if k := osMgr.Get("Hello"); k != "World" {
		t.Errorf("expected %v, got %v", "World", k)
	}

	if k, ok := osMgr.Lookup("Hello"); k != "World" || !ok {
		t.Errorf("expected a manager to be initialized empty for lookups, got %v - %v", k, ok)
	}
}
