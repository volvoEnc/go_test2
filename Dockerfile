FROM golang:latest

WORKDIR /go/src/app

COPY . .

EXPOSE 7020
