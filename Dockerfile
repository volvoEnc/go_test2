FROM golang:latest

WORKDIR /go/src/app

COPY . .

RUN go run main.go
