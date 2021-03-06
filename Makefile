build:
	go build -o ./price_rate_ticker -v

test:
	go test -v -short -race ./...

lint:
	golangci-lint run --verbose

fix-lint:
	golangci-lint run --verbose --fix

mock:
	# Remove old mockery files
	rm -rf mocks/*
	# Generate mock interfaces
	mockery --all --keeptree

check: build test lint

fix: fix-lint
