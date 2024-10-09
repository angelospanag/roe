.PHONY: run update-deps update-all-deps build generate-client

run:
	go run main.go

update-deps:
	go get -u ./... && go mod tidy

build:
	go build -o roe

update-all-deps:
	go get -u ./... && go mod tidy && cd frontend && pnpm update

generate-client:
	cd frontend && pnpm generate-client
