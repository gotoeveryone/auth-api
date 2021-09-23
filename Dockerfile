FROM golang:1.16-alpine

WORKDIR /go/src/github.com/auth-api

RUN apk add gcc g++

COPY go.mod go.sum ./
RUN go mod download

CMD ["go", "run", "main.go"]
