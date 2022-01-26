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
$ docker-compose up
```

## Test

```
$ docker-compose exec api go test
```

## Build

```
$ docker-compose exec api go build main.go
```

## Endpoint

|Endpoint|Method|Required token|Description|
|:--|:--|:--|:--|
|/|GET||Get application states|
|/users|GET|O|Get account data|
|/users|POST||Execute account registration (and issue a temporary password)|
|/activate|POST||Enabled account (and change password)|
|/auth|POST|O|Execute authentication|
|/deauth|DELETE|O|Execute deauthentication|

## Usage

- Example with `/auth`

### Request

- HTTP Header
  - Content-Type: application/json

```json
{
  "account": "test",
  "password": "password"
}
```

### Response

```
{
  "id": 1,
  "accessToken": "[Access Token]",
  "environment": "debug"
}
```

- Example at `/users`

### Request

- HTTP Header
  - Content-Type: application/json
  - Authorization: Bearer [Access Token]

### Response

```json
{
  "id": 1,
  "account": "test",
  "name": "Test User",
  "sex": "Male",
  "mailAddress": "test@example.com",
  "role": "Administrator"
}
```
