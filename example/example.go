package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/go-franky/keys"
	awsMgr "github.com/go-franky/keys/aws"
	"github.com/go-franky/keys/internal/aws"
	"github.com/go-franky/keys/os"
)

func main() {
	var secretID = flag.String("secret-id", "", "AWS Secret Manager ID")
	var region = flag.String("region", "us-west-1", "AWS Region for secrets")
	flag.Parse()

	cfg := aws.MustNewConfig(context.Background(), config.WithRegion(*region))
	sm := secretsmanager.NewFromConfig(
		cfg,
	)

	keyManager := keys.MultiManager(os.NewFromOS(), awsMgr.NewAWSKeyManager(*secretID, sm))
	fmt.Println(keyManager.Get("PASSWORD"))
}
