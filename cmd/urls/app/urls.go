package app

import (
	"sync"
	"time"
)

type urlsMap map[string]time.Duration

type minMax struct {
	mx     sync.RWMutex
	minUrl string
	min    time.Duration
	maxUrl string
	max    time.Duration
}

type urls struct {
	mx sync.RWMutex
	m  urlsMap

	minMax minMax
}

func (u *urls) Load(key string) (time.Duration, bool) {
	u.mx.RLock()
	defer u.mx.RUnlock()
	val, ok := u.m[key]
	return val, ok
}

func (u *urls) Store(key string, value time.Duration) {
	u.mx.Lock()
	defer u.mx.Unlock()
	u.m[key] = value
}

func (u *urls) FindMinMax() {
	u.mx.RLock()
	defer u.mx.RUnlock()

	for url, t := range u.m {
		if t < u.minMax.min && t != -1 {
			u.minMax.SetMin(url, t)
		} else if t > u.minMax.max {
			u.minMax.SetMax(url, t)
		}
	}
}

func (m *minMax) GetMin() (string, time.Duration) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.minUrl, m.min
}

func (m *minMax) SetMin(url string, t time.Duration) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.min = t
	m.minUrl = url
}

func (m *minMax) GetMax() (string, time.Duration) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.maxUrl, m.max
}

func (m *minMax) SetMax(url string, t time.Duration) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.max = t
	m.maxUrl = url
}
