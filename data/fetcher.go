package data

type Fetcher interface {
	Run() error
	Stop() error
}
