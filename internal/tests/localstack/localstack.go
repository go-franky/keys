package localstack

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	internalAWS "github.com/go-franky/keys/internal/aws"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func Setup() (func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return func() {}, fmt.Errorf("could not connect to docker: %w", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "localstack/localstack",
		Tag:        "latest",
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return func() {}, fmt.Errorf("could not start resource: %w", err)
	}

	close := func() {
		if err := pool.Purge(resource); err != nil {
			panic(err)
		}
	}

	hostAndPort := resource.GetHostPort("4566/tcp")
	os.Setenv("KEYS_AWS_ENDPOINT", fmt.Sprintf("http://%s", hostAndPort))

	resource.Expire(60)
	pool.MaxWait = 1 * time.Minute

	if err := pool.Retry(func() error {
		cfg := internalAWS.MustNewConfig(
			context.Background(),
			config.WithRetryer(func() aws.Retryer { return aws.NopRetryer{} }),
		)
		s := sts.NewFromConfig(cfg)
		_, err = s.GetCallerIdentity(context.Background(), nil)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return close, err
	}

	return close, nil
}
