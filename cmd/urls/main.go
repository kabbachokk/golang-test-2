package main

import (
	"context"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
	"unicode"

	"github.com/joho/godotenv"
	"work.test/cmd/urls/app"
)

func main() {
	godotenv.Load()

	var (
		username = flag.String("u", os.Getenv("AUTH_USERNAME"), "username")
		password = flag.String("p", os.Getenv("AUTH_PASSWORD"), "password")
		file     = flag.String("f", "urls.txt", "urls")
	)

	flag.Parse()
	if *username == "" {
		log.Fatal("не указано имя пользователя")
	}
	if *password == "" {
		log.Fatal("не указан пароль")
	}

	content, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("невозможно прочитать файл %v", err)
	}
	urls := strings.Split(string(content), "\n")
	for i, url := range urls {
		clean := strings.Map(func(r rune) rune {
			if unicode.IsPrint(r) {
				return r
			}
			return -1
		}, url)
		urls[i] = clean
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	app := app.NewApp(*username, *password, urls, client)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/min", app.MinHandler)
	mux.HandleFunc("/api/max", app.MaxHandler)
	mux.HandleFunc("/api", app.UrlHandler)
	mux.HandleFunc("/api/stats", app.AuthMiddleware(app.StatsHandler))

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	ticker := time.NewTicker(time.Minute)
	go app.UpdateUrls()
	go func() {
		for range ticker.C {
			app.UpdateUrls()
		}
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ticker.Stop()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
