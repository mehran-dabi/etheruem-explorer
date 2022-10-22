package api

import (
	"context"
	_ "embed"
	"energi-challenge/config"
	"energi-challenge/console/indexer/repository"
	"energi-challenge/infrastructure/sqlite"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/twiny/ratelimit"
)

type Indexer struct {
	wg          *sync.WaitGroup
	conf        *config.Configs
	srv         *http.Server
	limiter     *ratelimit.Limiter
	client      *ethclient.Client
	jobs        chan int64 // queue of block ids to scan
	subscribed  bool       // used to only subscribe to `client.SubscribeNewHead` once.
	latestBlock int64
	events      chan *types.Header
	store       sqlite.ISqlite
	repository  repository.IIndexerRepository
	ctx         context.Context
	done        context.CancelFunc
}

func NewIndexer() (*Indexer, error) {
	conf := config.Init()

	// db
	store, err := sqlite.NewSQLiteDB(conf.Store.Path)
	if err != nil {
		return nil, err
	}

	if err := store.Ping(); err != nil {
		return nil, err
	}

	if err := store.Migrate("up"); err != nil {
		return nil, err
	}

	client, err := ethclient.Dial(conf.Indexer.Endpoint)
	if err != nil {
		return nil, err
	}

	// get current latest block
	latest, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// api ctx
	ctx, done := context.WithCancel(context.Background())

	// api
	idx := &Indexer{
		wg:          &sync.WaitGroup{},
		conf:        conf,
		limiter:     ratelimit.NewLimiter(conf.Indexer.Limiter.Rate, conf.Indexer.Limiter.Duration),
		client:      client,
		jobs:        make(chan int64, conf.Indexer.ChanSize),
		subscribed:  false,
		latestBlock: latest.Number().Int64(),
		events:      make(chan *types.Header, conf.Indexer.ChanSize),
		store:       store,
		repository:  repository.NewIndexerRepository(store.DB),
		ctx:         ctx,
		done:        done,
	}

	// start api
	idx.indexer()

	return idx, nil
}

// Start starts the indexer microservice
func (idx *Indexer) Start() error {
	// add routes
	server := idx.run(idx.conf.Indexer.Address)
	idx.srv = server

	return nil
}

// Shutdown shuts down the server
func (idx *Indexer) Shutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigs
	log.Println("shutting down ...")

	if err := idx.srv.Shutdown(context.TODO()); err != nil {
		log.Println(err)
	}

	idx.wg.Wait()

	close(idx.jobs)
	close(idx.events)

	log.Println("program exited")
	os.Exit(0)
}
