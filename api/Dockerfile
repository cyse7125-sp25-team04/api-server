FROM golang:1.24.0-alpine3.21
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 go build -o /go/bin/webapp -ldflags "-X main.version=1.0.0" ./cmd/main.go
EXPOSE 8080
ENTRYPOINT ["/go/bin/webapp"]