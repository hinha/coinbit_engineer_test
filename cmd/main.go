package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lovoo/goka"
	"golang.org/x/sync/errgroup"

	coin "github.com/hinha/coinbit_test"
	"github.com/hinha/coinbit_test/api"
	"github.com/hinha/coinbit_test/constanta"
	"github.com/hinha/coinbit_test/pb"
	"github.com/hinha/coinbit_test/provider"
)

var (
	service string
	broker  string
	wallet  *string
	remove  *bool
)

func init() {
	flag.StringVar(&service, "service", "", "running service ex: api, broker")
	flag.StringVar(&broker, "broker_name", "balance", "running broker")
	wallet = flag.String("wallet", "", "flagged as above-threshold deposit")
	remove = flag.Bool("remove", false, "remove entities of wallet_id")
	flag.Parse()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	grp, ctx := errgroup.WithContext(ctx)

	// provide Broker running
	providers := provider.NewBroker([]string{constanta.BROKER})
	if service == "api" {
		mux := api.New([]string{constanta.BROKER})
		defer mux.Close()

		srv := &http.Server{
			Handler:      mux,
			Addr:         fmt.Sprintf(":%s", constanta.PORT),
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,
		}
		log.Fatal(srv.ListenAndServe())
	} else if service == "broker" {
		grp.Go(func() error {
			return providers.Execute(ctx, coin.Group(broker))
		})
	} else {
		emitter, err := goka.NewEmitter([]string{constanta.BROKER}, coin.FlagWalletStream, new(coin.FlagEventEncoder))
		if err != nil {
			panic(err)
		}
		defer emitter.Finish()

		err = emitter.EmitSync(*wallet, &pb.FlagEvent{FlagRemoved: *remove})
		if err != nil {
			panic(err)
		}
		return
	}

	// Wait for SIGINT/SIGTERM
	waiter := make(chan os.Signal, 1)
	signal.Notify(waiter, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-waiter:
	case <-ctx.Done():
	}
	cancel()
	if err := grp.Wait(); err != nil {
		log.Println(err)
	}

	if err := grp.Wait(); err != nil {
		log.Println(err)
	}

	os.Exit(0)
}
