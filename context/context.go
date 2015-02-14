package context
import(
	"../scheduler"
	"../downloader"
	"../analyzer"
)

type Context struct {
	scheduler *scheduler.Scheduler
	downloaders []*downloader.Downloader
	analyzers []*analyzer.Analyzer
}

func New(scheduler *scheduler.Scheduler, downloaders []*downloader.Downloader, analyzers []*analyzer.Analyzer) *Context {
	context := &Context{scheduler:scheduler, downloaders:downloaders, analyzers:analyzers}
	return context
}

func (context *Context) Scheduler() *scheduler.Scheduler {
	return context.scheduler
}

func (context *Context) Downloaders() []*downloader.Downloader {
	return context.downloaders
}

func (context *Context) Analyzers() []*analyzer.Analyzer {
	return context.analyzers
}