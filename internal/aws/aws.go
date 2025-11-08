package aws

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

func NewConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	if os.Getenv("KEYS_AWS_DEBUG") == "true" {
		optFns = append(optFns, config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody))
		optFns = append(optFns, config.WithRetryer(func() aws.Retryer { return aws.NopRetryer{} }))
	}
	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("unable to load SDK config, %w", err)
	}
	return cfg, nil
}

func MustNewConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) aws.Config {
	cfg, err := NewConfig(ctx, optFns...)
	if err != nil {
		panic(err)
	}
	return cfg
}

// customResolver implements secretsmanager.EndpointResolverV2
type customResolver[T any] struct {
	endpoint string
}

func (r *customResolver[T]) ResolveEndpoint(ctx context.Context, params T) (smithyendpoints.Endpoint, error) {
	if r.endpoint == "" {
		return smithyendpoints.Endpoint{}, fmt.Errorf("endpoint not set")
	}
	u, err := url.Parse(r.endpoint)
	if err != nil {
		return smithyendpoints.Endpoint{}, fmt.Errorf("failed to parse endpoint: %w", err)
	}
	return smithyendpoints.Endpoint{
		URI: *u,
	}, nil
}

func NewCustomResolver[T any](endpoint string) *customResolver[T] {
	return &customResolver[T]{
		endpoint: endpoint,
	}
}
