# General authentication API

Execute authentication and authorization with JSON Request and Response.

## Getting Started

### Prerequisites

Installed the following program.

- Golang 1.8+

## Installing

- Get `dep` binary (`dep` is Golang's package manager)

```
$ go get -u github.com/golang/dep/cmd/dep
$
$ cd <repository_root>
$ dep ensure
```

## Build & Run

```
$ GOOS=<target_os> GOARCH=<target_arch> go build -o general-api
$ ./general-api
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
