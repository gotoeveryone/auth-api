FROM golang:1.16-alpine

ENV GOPATH /go

WORKDIR /go/src/app

RUN apk add gcc g++

COPY go.mod go.sum ./
RUN go mod download && \
  go get github.com/cosmtrek/air

CMD ["air", "-c", ".air.toml"]
