package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/go-franky/keys"
	awsMgr "github.com/go-franky/keys/aws"
	"github.com/go-franky/keys/os"
)

func main() {
	var secretID = flag.String("secret-id", "", "AWS Secret Manager ID")
	var region = flag.String("region", "us-west-1", "AWS Region for secrets")
	flag.Parse()

	sess, err := session.NewSession(&aws.Config{Region: region})
	if err != nil {
		log.Fatalf("could not create a session: %v", err)
	}
	sm := secretsmanager.New(sess)

	keyManager := keys.Combine(os.NewFromOS(), awsMgr.NewAWSKeyManager(*secretID, sm))
	fmt.Println(keyManager.Get("PASSWORD"))
}
