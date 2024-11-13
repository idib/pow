# Stage for build app
FROM golang:1.23.2 AS builder

WORKDIR /app

COPY . ./

RUN go mod download
RUN go build -o server ./cmd/server
ENTRYPOINT  ["./server"]
