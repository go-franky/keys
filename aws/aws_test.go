//go:build localstack

package aws_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	awsMgr "github.com/go-franky/keys/aws"
	"github.com/go-franky/keys/internal/aws"
	"github.com/go-franky/keys/internal/tests/localstack"
)

func TestKeys(t *testing.T) {
	close, err := localstack.Setup()
	if err != nil {
		t.Fatal(err)
	}
	defer close()

	secretID := "anything"

	s := secretsmanager.NewFromConfig(
		aws.MustNewConfig(context.Background()),
		secretsmanager.WithEndpointResolverV2(
			aws.NewCustomResolver[secretsmanager.EndpointParameters](os.Getenv("KEYS_AWS_ENDPOINT")),
		),
	)

	data := struct {
		Hello string `json:"HELLO"`
	}{Hello: "World!"}

	val, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	text := string(val)
	if _, err := s.CreateSecret(context.Background(), &secretsmanager.CreateSecretInput{
		Name:         &secretID,
		SecretString: &text,
	}); err != nil {
		t.Fatalf("could not create secret: %v", err)
	}

	mgr := awsMgr.NewAWSKeyManager(secretID, s)

	if res := mgr.Get("HELLO"); res != data.Hello {
		t.Fatalf("expteded %v, got %v", data.Hello, res)
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
