package store

type Store interface {
	GetFullUrl(shortUrl string) string
	GetShortUrl(fullUrl string) string
	SetUrl(fullUrl, shortUrl string) bool
	CountUrls() int
	FlushUrls()
	RecordClick(fullUrl string)
	GetUrlHits(fullUrl string) int
	GetAllHits() map[string]int
	FlushHits()
}
