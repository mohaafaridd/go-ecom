# E-commerce App

This app is built following [Complete Backend API in Golang (JWT, MySQL & Tests)](https://www.youtube.com/watch?v=7VLmLOiQ3ck) tutorial on Youtube

### Tools
* Go 1.24.2
* [Database Migration CLI](https://github.com/golang-migrate/migrate/tree/v4.18.2/cmd/migrate)

### Booting up steps
1. `make migrate-up`
2. `make run`

### Features used

1. HTTP Gorilla Mux server
1. MySQL Database
1. [JWT](https://pkg.go.dev/github.com/golang-jwt/jwt/v5@v5.2.2) for authentication
1. [bcrypt](https://pkg.go.dev/golang.org/x/crypto@v0.37.0/bcrypt) for hashing
1. [Input validation library](https://pkg.go.dev/github.com/go-playground/validator@v9.31.0+incompatible)
