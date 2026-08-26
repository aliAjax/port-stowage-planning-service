.PHONY: build test vet run smoke
build:
	go build ./...
test:
	go test ./...
vet:
	go vet ./...
run:
	go run ./cmd/planner
smoke:
	./scripts/smoke.sh
