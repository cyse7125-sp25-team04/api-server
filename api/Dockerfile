FROM golang:1.24.0-alpine3.21
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o /go/bin/webapp ./cmd/main.go
EXPOSE 8080
ENTRYPOINT ["/go/bin/webapp"]