package shorten

import (
	store "./store"
	strategy "./strategy"
	"errors"
	"sync"
)

type Shortener struct {
	store       store.Store
	generateKey strategy.GenerateKey
	lock        sync.Mutex
}

type Options struct {
	Store       store.Store
	GenerateKey strategy.GenerateKey
}

const (
	NumRetries = 100
)

func New(o *Options) *Shortener {
	s := o.Store
	if s == nil {
		s = store.NewRedis(&store.RedisOptions{})
	}
	genkey := o.GenerateKey
	if genkey == nil {
		genkey = strategy.DefaultRandomKey()
	}

	return &Shortener{
		store:       s,
		generateKey: genkey,
	}
}

func (s *Shortener) Shorten(url string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	key := s.store.GetShortUrl(url)
	if key != "" {
		return key, nil
	}

	for i := 0; i < NumRetries; i++ {
		key = s.generateKey(s.store)
		ok := s.store.SetUrl(url, key)
		if ok {
			return key, nil
		}
	}
	return "", errors.New("Can't generate unique key after " + string(NumRetries) + " retries")
}

func (s *Shortener) Expand(shortUrl string) string {
	return s.store.GetFullUrl(shortUrl)
}

func (s *Shortener) RecordClick(url string) {
	s.store.RecordClick(url)
}
