test:
	go test -count 1 -cover -v -race ./...

vet:
	go vet ./...

fmt:
	test -z $$(gofmt -l .) # This will return non-0 if unsuccessful  run `go fmt ./...` to fix
