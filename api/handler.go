package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lovoo/goka"

	coin "github.com/hinha/coinbit_test"
	"github.com/hinha/coinbit_test/constanta"
	"github.com/hinha/coinbit_test/provider"
)

// HTTPHandler is a http.Handler for application.
type HTTPHandler struct {
	router  *mux.Router
	closers []func() error
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// Close closes connections client.
func (h *HTTPHandler) Close() error {
	for _, f := range h.closers {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// New creates a HTTPHandler with the given brokers.
func New(brokers []string) *HTTPHandler {

	view, err := goka.NewView(brokers, provider.BalanceGroupTable, new(coin.DepositListEncoder))
	if err != nil {
		panic(err)
	}
	go view.Run(context.Background())

	flaggerView, err := goka.NewView(brokers, provider.FlagGroupTable, new(coin.FlagValueEncoder))
	if err != nil {
		panic(err)
	}
	go flaggerView.Run(context.Background())

	emitter, err := goka.NewEmitter(brokers, coin.DepositStream, new(coin.DepositEncoder))
	if err != nil {
		panic(err)
	}

	return &HTTPHandler{
		router:  muxRouter(emitter, coin.DepositStream, view, flaggerView),
		closers: []func() error{emitter.Finish},
	}
}

func muxRouter(emitter *goka.Emitter, stream goka.Stream, balanceView, flagView *goka.View) *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	// Endpoints.
	api.HandleFunc("/deposit", newDepositHandlerFunc(emitter, stream)).Methods(http.MethodPost)
	api.HandleFunc("/check/{wallet_id}", newCheckHandlerFunc(emitter, balanceView, flagView)).Methods(http.MethodGet)

	log.Printf("Listen api port %s", constanta.PORT)
	return router
}
