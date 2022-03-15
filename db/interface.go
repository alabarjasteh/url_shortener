package db

type Interface interface {
	Load(shortLink string) (string, error)
	Save(shortLink, originalLink string) error
}
