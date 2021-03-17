package main

import (
	"net/http"
	"os"
	"sync"
	"time"

	"proxy/internal/app"
	"proxy/internal/infrastructure/storage"
	"proxy/internal/interfaces/web"
	"proxy/internal/interfaces/web/server"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

// const (
// 	prefix = "proxy"
// )

func main() {
	SetupLogger(os.Stdout, "info")

	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:                 "db",
			Port:                 5432,
			User:                 "user",
			Database:             "db",
			Password:             "password",
			PreferSimpleProtocol: true,
		},
		MaxConnections: 10,
		AcquireTimeout: time.Duration(3 * time.Second),
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("Cannnot connect to db", err)
		os.Exit(1)
	}
	defer conn.Close()

	db, err := storage.NewDB(conn)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("Cannot create infrastructure layer", err)
		os.Exit(1)
	}

	app := app.NewApp(db)

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		sniffercfg := &server.Config{
			Addr:              "0.0.0.0:8000",
			ReadTimeout:       time.Duration(3 * time.Second),
			ReadHeaderTimeout: time.Duration(3 * time.Second),
			WriteTimeout:      time.Duration(3 * time.Second),
			IdleTimeout:       time.Duration(3 * time.Second),
		}

		sniffer := server.NewServer(http.HandlerFunc(web.NewSniffer(app, sniffercfg).Recording), sniffercfg)

		err := sniffer.ListenAndServe()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "main",
				"func": "main",
			}).Fatalln("Sniffer server error", err)
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		admincfg := &server.Config{
			Addr:              "0.0.0.0:8080",
			ReadTimeout:       time.Duration(3 * time.Second),
			ReadHeaderTimeout: time.Duration(3 * time.Second),
			WriteTimeout:      time.Duration(3 * time.Second),
			IdleTimeout:       time.Duration(3 * time.Second),
		}

		admin := server.NewServer(web.NewAdmin(app, admincfg).Init(), admincfg)
		err := admin.ListenAndServe()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"pack": "main",
				"func": "main",
			}).Fatalln("Admin server error", err)
		}
	}(&wg)

	wg.Wait()
}
