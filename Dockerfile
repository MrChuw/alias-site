FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

COPY templates/ /app/templates/

RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o main

ENTRYPOINT ["/app/main"]

EXPOSE 80


