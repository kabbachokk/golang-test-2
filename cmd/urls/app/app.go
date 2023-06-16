package app

import (
	"net/http"
	"time"
)

type App struct {
	auth struct {
		username string
		password string
	}
	urls   urls
	client http.Client
}

func NewApp(username string, password string, urls []string, client http.Client) *App {
	app := new(App)
	app.auth.username = username
	app.auth.password = password
	app.urls.m = fillMap(urls)
	app.urls.minMax.max = 0
	app.urls.minMax.maxUrl = ""
	app.urls.minMax.min = time.Duration(time.Hour)
	app.urls.minMax.minUrl = ""
	app.client = client

	return app
}

func fillMap(urls []string) urlsMap {
	var m urlsMap
	m = make(urlsMap)
	for _, s := range urls {
		m[s] = 0
	}
	return m
}
