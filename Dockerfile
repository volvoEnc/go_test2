FROM golang:latest

WORKDIR /go/src/app

RUN go run main.go
