package app

import (
	"encoding/json"
	"net/http"
	"time"
)

func (app *App) StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	app.urls.mx.RLock()
	jsonRes, err := json.Marshal(app.urls.m)
	app.urls.mx.RUnlock()

	if err != nil {
		http.Error(w, "внутренняя ошибка", 500)
	}

	w.Write(jsonRes)
	return
}

func (app *App) MinHandler(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]time.Duration)
	url, t := app.urls.minMax.GetMin()
	res[url] = t

	w.Header().Set("Content-Type", "application/json")
	jsonRes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "внутренняя ошибка", 500)
	}
	w.Write(jsonRes)
	return
}

func (app *App) MaxHandler(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]time.Duration)
	url, t := app.urls.minMax.GetMax()
	res[url] = t

	w.Header().Set("Content-Type", "application/json")
	jsonRes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "внутренняя ошибка", 500)
	}
	w.Write(jsonRes)
	return
}

func (app *App) UrlHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	res := make(map[string]time.Duration)
	t, ok := app.urls.Load(url)
	if !ok {
		http.Error(w, "указанный url не найден", 404)
	}
	res[url] = t

	w.Header().Set("Content-Type", "application/json")
	jsonRes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "внутренняя ошибка", 500)
	}
	w.Write(jsonRes)
	return
}
