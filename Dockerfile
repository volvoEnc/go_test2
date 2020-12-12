FROM golang:latest

WORKDIR /go/src/app

COPY . .

EXPOSE 7020

RUN go run main.go
