FROM golang:1.22 AS base

FROM base as dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

WORKDIR /app

RUN GOOS=linux go build -o /dealer-status ./cmd

WORKDIR /

ENTRYPOINT [ "/dealer-status" ]