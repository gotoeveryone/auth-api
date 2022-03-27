# General authentication API

![Build Status](https://github.com/gotoeveryone/auth-api/workflows/Build/badge.svg)
[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/gotoeveryone/myrecipe/blob/master/LICENSE)
[![GitHub version](https://badge.fury.io/gh/gotoeveryone%2Fauth-api.svg)](https://badge.fury.io/gh/gotoeveryone%2Fauth-api)

Execute authentication and authorization with JSON Request and Response.

## Getting Started

### Prerequisites

Installed the following program.

- Docker

## Run

```
$ cp .env.example .env # with editing
$ docker compose up
```

## Format Check

```
$ docker compose exec api go vet -v ./...
```

## Test

```
$ docker compose exec api go test -v ./...
```

## Build

```
$ docker compose exec api go build main.go
```

## Swagger UI

- /swagger/index.html
