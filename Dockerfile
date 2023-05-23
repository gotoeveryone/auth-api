FROM golang:1.19-alpine as development

ENV LANG C.UTF-8
ENV APP_ROOT /var/app

# hadolint ignore=DL3018
RUN apk add gcc g++ --no-cache

WORKDIR ${APP_ROOT}
COPY go.mod go.sum ./

ENV DOCKERIZE_VERSION v0.6.1
RUN wget -q https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
  && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
  && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

RUN go install github.com/cosmtrek/air@v1.29.0 && \
  go install github.com/swaggo/swag/cmd/swag@v1.8.12 && \
  go install honnef.co/go/tools/cmd/staticcheck@2022.1.2

# uncomment if use sql-migrate run migration instead of gorm
# RUN go install github.com/rubenv/sql-migrate/...@v1.1.1

RUN go mod download

CMD ["air", "-c", ".air.toml"]

FROM golang:1.19-alpine as builder

ENV LANG C.UTF-8
ENV APP_ROOT /var/app
ENV GIN_MODE release

# hadolint ignore=DL3018
RUN apk add gcc g++ --no-cache

WORKDIR ${APP_ROOT}
COPY ./ ${APP_ROOT}

RUN go mod download && \
  go build -o auth-api ${APP_ROOT}/app/main.go

FROM golang:1.19-alpine as production

ENV LANG C.UTF-8
ENV APP_ROOT /var/app
ENV GIN_MODE release

WORKDIR ${APP_ROOT}
COPY --from=builder ${APP_ROOT}/auth-api ${APP_ROOT}

CMD ["${APP_ROOT}/auth-api"]
