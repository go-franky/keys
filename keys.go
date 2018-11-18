// Package keys offers a small helper for managing
// multiple set of keys / values of the type string string
// It allows for keys to be searched from multiple source
package keys

// Getter retrieves the value of the variable named by the key.
// It returns the value, which will be empty if the variable is not present.
// To distinguish between an empty value and an unset value, use Lookup.
type Getter interface {
	Get(k string) string
}

// Lookuper retrieves the value of the variable named by the key.
// If the variable is present in the environment the value (which may be empty)
// is returned and the boolean is true. Otherwise the returned value will be empty
// and the boolean will be false.
type Lookuper interface {
	Lookup(k string) (string, bool)
}

// Setter sets the value of the variable named by the key.
// It returns an error, if any.
type Setter interface {
	Set(k, v string) error
}

// Manager is the inter interface to a key manager
type Manager interface {
	Getter
	Lookuper
	Setter
}

// KeyManager is a basic implementation in memory
// of a Manager
type KeyManager struct {
	localData map[string]string
}

// Lookup see interface definition
func (km *KeyManager) Lookup(key string) (string, bool) {
	k, v := km.localData[key]
	return k, v
}

// Get see interface definition
func (km *KeyManager) Get(key string) string {
	k, _ := km.Lookup(key)
	return k
}

// Set see interface definition
func (km *KeyManager) Set(key, value string) error {
	km.localData[key] = value
	return nil
}

// NewKeyManager gives a basic key manager that will
// store values in memory
func NewKeyManager() Manager {
	return &KeyManager{
		localData: make(map[string]string),
	}
}

type cbGet struct {
	mgr []Getter
}

func (km *cbGet) Get(key string) string {
	for _, keym := range km.mgr {
		if k := keym.Get(key); k != "" {
			return k
		}
	}
	return ""
}

type cbLook struct {
	mgr []Lookuper
}

func (km *cbLook) Lookup(key string) (string, bool) {
	for _, keym := range km.mgr {
		if k, ok := keym.Lookup(key); ok {
			return k, ok
		}
	}
	return "", false
}

type combinedManager struct {
	localData map[string]string
	mgr       []Manager
}

func (km *combinedManager) Get(key string) string {
	k, _ := km.Lookup(key)
	return k
}

func (km *combinedManager) Lookup(key string) (string, bool) {
	if k, ok := km.localData[key]; ok {
		return k, ok
	}
	for _, keym := range km.mgr {
		if k, ok := keym.Lookup(key); ok {
			return k, ok
		}
	}
	return "", false
}

func (km *combinedManager) Set(k, v string) error {
	km.localData[k] = v
	return nil
}

// MultiGetter takes multiple lookups and combines them to find
// the value of the first one that is not a blank string
func MultiGetter(km ...Getter) Getter {
	res := &cbGet{}
	for _, k := range km {
		res.mgr = append(res.mgr, k)
	}
	return res
}

// MultiLookuper takes multiple lookups and combines them to find
// the value of any
func MultiLookuper(km ...Lookuper) Lookuper {
	res := &cbLook{}
	for _, k := range km {
		res.mgr = append(res.mgr, k)
	}
	return res
}

// MultiManager takes multiple managers and combines them to find
// the value of any
func MultiManager(km ...Manager) Manager {
	res := &combinedManager{
		localData: make(map[string]string),
	}
	for _, k := range km {
		res.mgr = append(res.mgr, k)
	}
	return res
}
