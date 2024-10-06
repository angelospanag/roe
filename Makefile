.PHONY: run update-deps build

run:
	go run main.go

update-deps:
	go get ./... && go mod tidy

build:
	go build -o roe
