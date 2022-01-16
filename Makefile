GOTEST=go test -v -count=1 -race -cover

test:
	$(GOTEST) -race ./...

test_with_localstack:
	$(GOTEST) -tags=localstack \
		github.com/go-franky/keys/aws

vet:
	go vet ./...

fmt:
	test -z $$(gofmt -l .) # This will return non-0 if unsuccessful  run `go fmt ./...` to fix
