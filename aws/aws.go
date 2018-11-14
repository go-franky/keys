package aws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/go-franky/keys"
)

type awsKeyManager struct {
	once      sync.Once
	sm        *secretsmanager.SecretsManager
	secretID  string
	localData map[string]string
}

func (km *awsKeyManager) Lookup(key string) (string, bool) {
	km.once.Do(km.getData)
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

func (km *awsKeyManager) getData() {
	secretValue, err := km.sm.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: aws.String(km.secretID)})
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

// NewAWSKeyManager is a concrete implementation of keys.Manager which interacts with
// AWS Secrets Manager.
func NewAWSKeyManager(secretID string, sm *secretsmanager.SecretsManager) keys.Manager {
	return &awsKeyManager{
		sm:        sm,
		secretID:  secretID,
		localData: make(map[string]string),
	}
}
