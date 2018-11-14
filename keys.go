// Package keys offers a small helper for managing
// multiple set of keys / values of the type string string
// It allows for keys to be searched from multiple source
package keys

// Getter is the interface for accessing a key
type Getter interface {
	Get(k string) string
}

// Lookuper is the interface for accessing a key
// but returning a bool when no values are present
// This can help differentiate no value from an empty value
type Lookuper interface {
	Lookup(k string) (string, bool)
}

// Manager is the inter interface to a key manager
type Manager interface {
	Getter
	Lookuper
	Set(k, v string) error
}

// Combine takes multiple managers and combines them to find
// the value of any
func Combine(km ...Manager) Manager {
	res := &combinedManager{
		localData: make(map[string]string),
	}
	for _, k := range km {
		res.mgr = append(res.mgr, k)
	}
	return res
}

// KeyManager is a basic implementation in memory
// of a Manager
type KeyManager struct {
	localData map[string]string
}

func (km *KeyManager) Lookup(key string) (string, bool) {
	k, v := km.localData[key]
	return k, v
}

func (km *KeyManager) Get(key string) string {
	k, _ := km.Lookup(key)
	return k
}

func (km *KeyManager) Set(key, value string) error {
	km.localData[key] = value
	return nil
}

// NewKeyManager gives a basic key manager
func NewKeyManager() Manager {
	return &KeyManager{
		localData: make(map[string]string),
	}
}

type combinedManager struct {
	localData map[string]string
	mgr       []Manager
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

func (km *combinedManager) Get(key string) string {
	k, _ := km.Lookup(key)
	return k
}
