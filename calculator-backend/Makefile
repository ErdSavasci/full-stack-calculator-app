test:
	go clean -testcache
	go list -f '{{.Dir}}/...' -m | xargs go test -v -cover

coverage:
	go clean -testcache
	go list -f '{{.Dir}}/...' -m | xargs go test -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out