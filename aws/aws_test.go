package aws_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	awsMgr "github.com/go-franky/keys/aws"
)

type fakeSecretValuer struct{}

func (f *fakeSecretValuer) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return &secretsmanager.GetSecretValueOutput{
		SecretString: aws.String(`{"Hello":"World"}`),
	}, nil
}

func TestKeys(t *testing.T) {
	mgr := awsMgr.NewAWSKeyManager("anything", &fakeSecretValuer{})

	if res := mgr.Get("Hello"); res != "World" {
		t.Fatalf("expteded %v, got %v", "Hello", res)
	}

	if res := mgr.Get("NAME"); res != "" {
		t.Fatalf("expected %v to be empty", res)
	}

	if err := mgr.Set("NAME", "franky"); err != nil {
		t.Fatal(err)
	}

	if res := mgr.Get("NAME"); res != "franky" {
		t.Fatalf("expected %v, got %v", "franky", res)
	}
}
