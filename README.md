## Chirpy

Chirpy is a http server project built by using Go and PostgreSQL 

## Goal

Chirpy goal is to provide API functionality close to social network like Twitter

## Installation
1.Install go

3.Install goose
```go install github.com/pressly/goose/v3/cmd/goose@latest```

4.Clone repo

5.Run goose up migrations with your postgres db uri

6.Configure .env variables, for example:
DB_URL=""
PLATFORM=""
SECRETKEY=""
POLKA_KEY=""

7.Run the server by using next command in main directory:
```go run .```

## Usage

After running the server, you can access it at http://localhost:8080
All endpoints are described in API documentation 

## Other tools

I also used sqlc instead of database package, GORM, etc.
To Install sqlc: 
```go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest```
