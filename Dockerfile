FROM golang:1.16-alpine as development

ENV LANG C.UTF-8
ENV APP_ROOT /var/app

RUN apk add gcc g++

WORKDIR ${APP_ROOT}
COPY go.mod go.sum ./

RUN go mod download && \
  go get -u github.com/cosmtrek/air && \
  go get -u github.com/swaggo/swag/cmd/swag

CMD ["air", "-c", ".air.toml"]

FROM golang:1.16-alpine as builder

ENV LANG C.UTF-8
ENV APP_ROOT /var/app
ENV GIN_MODE release

RUN apk add gcc g++

WORKDIR ${APP_ROOT}
COPY ./ ${APP_ROOT}

RUN go mod download && \
  go build -o auth-api ${APP_ROOT}/main.go

FROM golang:1.16-alpine as production

ENV LANG C.UTF-8
ENV APP_ROOT /var/app
ENV GIN_MODE release

WORKDIR ${APP_ROOT}
COPY --from=builder ${APP_ROOT}/auth-api ${APP_ROOT}

CMD ${APP_ROOT}/auth-api
