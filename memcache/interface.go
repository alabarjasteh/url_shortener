package memcache

type Interface interface {
	Get(shortLink string) (string, error)
	Set(shortLink, originalLink string) error
}
