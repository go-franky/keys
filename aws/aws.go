// Package AWS allows for storing AWS secretsmanager's secret
// into a local structure.
// Note: this assumes the value of the secret is in the format
// of a key value pair
package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/go-franky/keys"
)

type awsKeyManager struct {
	once      func() error
	localData map[string]string
}

func (km *awsKeyManager) Lookup(key string) (string, bool) {
	if err := km.once(); err != nil {
		log.Fatal(err)
	}

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

type getSecretValueer interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// NewAWSKeyManager is a concrete implementation of keys.Manager which interacts with
// AWS Secrets Manager.
func NewAWSKeyManager(secretID string, sm getSecretValueer) keys.Manager {
	km := &awsKeyManager{
		localData: make(map[string]string),
	}
	km.once = sync.OnceValue(km.loadData(secretID, sm))

	return km
}

func (km *awsKeyManager) loadData(secret string, sm getSecretValueer) func() error {
	return func() error {
		secretValue, err := sm.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{SecretId: &secret})
		if err != nil {
			return fmt.Errorf("could not get the secret value: %w", err)
		}
		var keys map[string]string
		if err := json.Unmarshal([]byte(*secretValue.SecretString), &keys); err != nil {
			return fmt.Errorf("could not unmarshal: %w", err)
		}
		for k, v := range keys {
			if err := km.Set(k, v); err != nil {
				return fmt.Errorf("could not set %v: %w", k, err)
			}
		}
		return nil
	}
}
