# go-ergo-handler

Ergonomic HTTP handlers builder for Go.

## About

This library can help you building robust type-safe HTTP-handlers from reusable parsers and middlewares. 

## Installation

```bash
go get github.com/nktknshn/go-ergo-handler
```

## Example

```go

```

## Usage

example project 
https://github.com/nktknshn/go-ergo-handler-example
https://github.com/nktknshn/go-ergo-handler-example/tree/master/internal/adapters/http_adapter/handlers

query param
https://github.com/nktknshn/go-ergo-handler-example/blob/master/internal/adapters/http_adapter/handlers/get_books/query_param_cursor.go

router param
https://github.com/nktknshn/go-ergo-handler-example/blob/master/internal/adapters/http_adapter/handlers/get_book/get_book.go

payload
https://github.com/nktknshn/go-ergo-handler-example/blob/master/internal/adapters/http_adapter/handlers/create_book/create_book.go

custom http code
https://github.com/nktknshn/go-ergo-handler-example/blob/master/internal/adapters/http_adapter/handlers/create_favorite_book/create_favorite_book.go

custom parser
https://github.com/nktknshn/go-ergo-handler-example/blob/master/internal/adapters/http_adapter/handlers/handler_admin_role_checker/handler_admin_role_checker.go

custom error handler
https://github.com/nktknshn/go-ergo-handler-example/blob/master/internal/adapters/http_adapter/handlers/handler_builder/handler_error.go
