package polyglot

type Downloader interface {
	Download(url string, dest string) (string, error)
}
