# go-ergo-handler

Ergonomic HTTP handlers builder for Go.

## About

This library can help you building robust type-safe HTTP-handlers from reusable middlewares. 

## Installation

```bash
go get github.com/nktknshn/go-ergo-handler
```

## Example

```go

type payloadType struct {
	Title string `json:"title"`
	Price int    `json:"price"`
}

func (p payloadType) Validate() error {
	if p.Title == "" {
		return errors.New("empty book title")
	}
	if p.Price <= 0 {
		return errors.New("invalid book price")
	}
	return nil
}

type paramBookIDType int

func (p paramBookIDType) Parse(ctx context.Context, v string) (paramBookIDType, error) {
	vint, err := strconv.Atoi(v)
	if err != nil {
		return 0, errors.New("invalid book id")
	}
	return paramBookIDType(vint), nil
}

func (p paramBookIDType) Validate() error {
	if p <= 0 {
		return errors.New("invalid book id")
	}
	return nil
}

var (
	paramBookID    = geh.RouterParamWithParser[paramBookIDType]("book_id", errors.New("book_id is required"))
	payloadBook    = geh.Payload[payloadType]()
	paramUnpublish = geh.QueryParamMaybe("unpublish", geh.IgnoreContext(strconv.ParseBool))
)

func makeHttpHandler(useCase interface {
	UpdateBook(ctx context.Context, bookID int, title string, price int, unpublish bool) error
}) http.Handler {
	var (
		builder   = geh.New()
		bookID    = paramBookID.Attach(builder)
		payload   = payloadBook.Attach(builder)
		unpublish = paramUnpublish.Attach(builder)
	)

	return builder.BuildHandlerWrapped(func(w http.ResponseWriter, r *http.Request) (any, error) {
		// all values are parsed and validated at this point
		bid := bookID.Get(r)
		pl := payload.Get(r)
		unpublish := unpublish.GetDefault(r, false)
		err := useCase.UpdateBook(r.Context(), int(bid), pl.Title, pl.Price, unpublish)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
}
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

authentication
https://github.com/nktknshn/go-ergo-handler-example/blob/master/internal/adapters/http_adapter/handlers/handlers_user_auth/user_auth_parser.go

custom parser and authorization
https://github.com/nktknshn/go-ergo-handler-example/blob/master/internal/adapters/http_adapter/handlers/handler_admin_role_checker/handler_admin_role_checker.go

custom error handler
https://github.com/nktknshn/go-ergo-handler-example/blob/master/internal/adapters/http_adapter/handlers/handler_builder/handler_error.go
