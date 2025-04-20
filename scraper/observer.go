package scraper

type Observer interface {
	StartProcess() error
	EndProcess() error
	Success(book string, msg string)
	Error(book string, msg string)
}
