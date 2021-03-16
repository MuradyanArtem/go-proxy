package main

import (
	"flag"
	"os"
	"sync"

	"proxy/internal/app"
	"proxy/internal/infrastructure/storage"
	"proxy/internal/interfaces/web"
	"proxy/internal/interfaces/web/server"

	"github.com/jackc/pgx"
)

const (
	prefix = "proxy"
)

func main() {
	Parse(flag.CommandLine, os.Args[1:])

	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:                 dbconf.Host,
			Port:                 dbconf.Port,
			User:                 dbconf.User,
			Database:             dbconf.Database,
			Password:             dbconf.Password,
			PreferSimpleProtocol: dbconf.PreferSimpleProtocol,
		},
		MaxConnections: dbconf.MaxConnections,
		AcquireTimeout: dbconf.AcquireTimeout,
	})
	if err != nil {
		// log
		os.Exit(1)
	}
	defer conn.Close()

	db, err := storage.NewDB(conn)
	if err != nil {
		// log
		os.Exit(1)
	}

	app := app.NewApp(db)

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		sniffer := server.NewServer(web.NewSniffer(app, sniffercfg).Init(), sniffercfg)
		err := sniffer.ListenAndServe()
		if err != nil {
			//log
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		admin := server.NewServer(web.NewAdmin(app, admincfg).Init(), admincfg)
		err := admin.ListenAndServe()
		if err != nil {
			//log
		}
	}(&wg)

	wg.Wait()
}
