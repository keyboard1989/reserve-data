package stat

type Fetcher interface {
	Run() error
	Stop() error
}
