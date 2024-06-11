run:
	go run ./cmd/interceptor/main.go
run-test-srv:
	go run ./cmd/testserver/main.go

.PHONY: run, run-test-srv
