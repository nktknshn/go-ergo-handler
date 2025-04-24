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
	ergo "github.com/nktknshn/go-ergo-handler"
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
	paramBookID    = ergo.RouterParamWithParser[paramBookIDType]("book_id")
	paramUnpublish = ergo.QueryParamMaybe("unpublish", ergo.IgnoreContext(strconv.ParseBool))
	payloadBook    = ergo.Payload[payloadType]()
)

func makeHttpHandler(useCase interface {
	UpdateBook(ctx context.Context, bookID int, title string, price int, unpublish bool) error
}) http.Handler {

	var (
		builder   = ergo.New()
		bookID    = paramBookID.Attach(builder)
		unpublish = paramUnpublish.Attach(builder)
		payload   = payloadBook.Attach(builder)
	)

	return builder.BuildHandlerWrapped(func(w http.ResponseWriter, r *http.Request) (any, error) {
		// all values are parsed and validated at this point
		bid := bookID.Get(r)
		pl := payload.Get(r)
		unp := unpublish.GetDefault(r, false)
		err := useCase.UpdateBook(r.Context(), int(bid), pl.Title, pl.Price, unp)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
}

type useCase struct{}

func (u useCase) UpdateBook(ctx context.Context, bookID int, title string, price int, unpublish bool) error {
	fmt.Println("bookID:", bookID, "title:", title, "price:", price, "unpublish:", unpublish)
	return nil
}

func main() {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("PUT", "/books/1?unpublish=true", strings.NewReader(`{"title": "The Great Gatsby", "price": 100}`))

	handler := makeHttpHandler(useCase{})
	router := mux.NewRouter()
	router.Handle("/books/{book_id:[0-9]+}", handler)

	router.ServeHTTP(w, r)

	/*
		bookID: 1 title: The Great Gatsby price: 100 unpublish: true
		200
		{"result":{}}
	*/

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	w = httptest.NewRecorder()
	r = httptest.NewRequest("PUT", "/books/1", strings.NewReader(`{"price": 100}`))
	router.ServeHTTP(w, r)

	/*
		400
		{"error":"empty book title"}
	*/

	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

}
