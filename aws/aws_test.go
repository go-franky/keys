package aws_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/go-franky/keys/aws"
)

type mockSecretsManagerClient struct {
	secretsmanageriface.SecretsManagerAPI
}

func (c *mockSecretsManagerClient) GetSecretValue(ip *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	value := `{"Hello":"World"}`
	return &secretsmanager.GetSecretValueOutput{
		SecretString: &value,
	}, nil
}

func TestData(t *testing.T) {
	km := aws.NewAWSKeyManager("test", &mockSecretsManagerClient{})

	if km.Get("Hello") != "World" {
		t.Errorf("expected %v, got %v", "World", km.Get("Hello"))
	}
}
