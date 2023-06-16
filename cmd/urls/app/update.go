package app

import (
	"log"
	"net/http"
	"sync"
	"time"
)

func (app *App) UpdateUrls() {
	const workers = 5

	wg := new(sync.WaitGroup)
	in := make(chan string, 10*workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range in {
				t, err := getResponseTime(url)
				if err != nil {
					log.Printf("ошибка запроса: %v", err)
				}
				app.urls.Store(url, t)
			}
		}()
	}

	for url := range app.urls.m {
		if url != "" {
			in <- url
		}
	}
	close(in)
	wg.Wait()

	app.urls.FindMinMax()
}

func getResponseTime(url string) (time.Duration, error) {
	time_start := time.Now()
	_, err := http.Get("http://" + url)
	if err != nil {
		return -1, err
	}

	return time.Since(time_start), nil
}
