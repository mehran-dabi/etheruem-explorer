package api

import (
	"context"
	_ "embed"
	"energi-challenge/config"
	"energi-challenge/console/rest/repository"
	"energi-challenge/infrastructure/sqlite"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type API struct {
	conf       *config.Configs
	server     *http.Server
	store      sqlite.ISqlite
	repository repository.IRestRepository
}

func NewAPI() (*API, error) {
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

	return &API{
		conf:       conf,
		store:      store,
		repository: repository.NewRestRepository(store.DB),
	}, nil
}

func (a *API) Start() error {
	// add routes
	server := a.Run(a.conf.Rest.Port)
	a.server = server

	return nil
}

func (a *API) Shutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigs
	log.Println("shutting down ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Println(err)
	}

	log.Println("program exited")
	os.Exit(0)
}
