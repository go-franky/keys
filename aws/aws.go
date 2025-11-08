// Package AWS allows for storing AWS secretsmanager's secret
// into a local structure.
// Note: this assumes the value of the secret is in the format
// of a key value pair
package aws

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/go-franky/keys"
)

type awsKeyManager struct {
	once      sync.Once
	sm        *secretsmanager.Client
	secretID  string
	localData map[string]string
}

func (km *awsKeyManager) Lookup(key string) (string, bool) {
	km.once.Do(km.getData(context.Background()))
	key, ok := km.localData[key]
	return key, ok
}

func (km *awsKeyManager) Get(key string) string {
	k, _ := km.Lookup(key)
	return k
}

func (km *awsKeyManager) Set(k, v string) error {
	km.localData[k] = v
	return nil
}

func (km *awsKeyManager) getData(ctx context.Context) func() {
	return func() {
		secretValue, err := km.sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &km.secretID})
		if err != nil {
			log.Fatalf("could not get the secret value: %v", err)
		}
		var keys map[string]string
		if err := json.Unmarshal([]byte(*secretValue.SecretString), &keys); err != nil {
			log.Fatalf("could not unmarshal: %v", err)
		}
		for k, v := range keys {
			km.Set(k, v)
		}
	}
}

// NewAWSKeyManager is a concrete implementation of keys.Manager which interacts with
// AWS Secrets Manager.
func NewAWSKeyManager(secretID string, sm *secretsmanager.Client) keys.Manager {
	return &awsKeyManager{
		sm:        sm,
		secretID:  secretID,
		localData: make(map[string]string),
	}
}
