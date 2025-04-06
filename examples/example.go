package examples

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	goergohandler "github.com/nktknshn/go-ergo-handler"
)

type payloadType struct {
	Title string `json:"title"`
	Price int    `json:"price"`
}

func (p payloadType) Validate() error {
	if p.Title == "" {
		return errors.New("invalid book title")
	}
	if p.Price <= 0 {
		return errors.New("invalid book price")
	}
	return nil
}

func makeHttpHandler(useCase interface {
	UpdateBook(ctx context.Context, bookID int, payload payloadType) error
}) http.Handler {
	var (
		paramBookID = goergohandler.NewRouterParam("book_id", func(ctx context.Context, v string) (int, error) {
			return strconv.Atoi(v)
		})
		payloadBook = goergohandler.NewPayloadWithValidation[payloadType]()

		builder = goergohandler.New()
		bookID  = paramBookID.Attach(builder)
		payload = payloadBook.Attach(builder)
		handler = builder.BuildHandlerWrapped(func(w http.ResponseWriter, r *http.Request) (any, error) {
			bid := bookID.GetRequest(r)
			pl := payload.GetRequest(r)
			err := useCase.UpdateBook(r.Context(), bid, pl)
			if err != nil {
				return nil, err
			}
			return nil, nil
		})
	)

	return handler
}
