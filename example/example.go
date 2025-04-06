package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	geh "github.com/nktknshn/go-ergo-handler"
)

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
	UpdateBook(ctx context.Context, bookID int, payload payloadType, unpublish bool) error
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
		err := useCase.UpdateBook(r.Context(), int(bid), pl, unpublish)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
}

type useCase struct{}

func (u useCase) UpdateBook(ctx context.Context, bookID int, payload payloadType, unpublish bool) error {
	fmt.Println("bookID:", bookID, "payload:", payload, "unpublish:", unpublish)
	return nil
}

func main() {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/books/1", strings.NewReader(`{"title": "The Great Gatsby", "price": 100}`))

	handler := makeHttpHandler(useCase{})
	router := mux.NewRouter()
	router.Handle("/books/{book_id:[0-9]+}", handler)

	router.ServeHTTP(w, r)

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())
}
