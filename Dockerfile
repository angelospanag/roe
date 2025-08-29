FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM scratch

COPY --from=builder /out/api /api

USER 65532:65532

EXPOSE 8000
ENTRYPOINT ["/api"]
