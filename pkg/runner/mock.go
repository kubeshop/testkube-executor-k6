package runner

import (
	"log"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
)

// MockFetcher implements the Mock version of the content fetcher from 	"github.com/kubeshop/testkube/pkg/executor/content"
type MockFetcher struct {
	FetchFn                     func(content *testkube.TestContent) (path string, err error)
	FetchStringFn               func(str string) (path string, err error)
	FetchURIFn                  func(uri string) (path string, err error)
	FetchGitDirFn               func(repo *testkube.Repository) (path string, err error)
	FetchGitFileFn              func(repo *testkube.Repository) (path string, err error)
	FetchGitFn                  func(repo *testkube.Repository) (path string, err error)
	FetchCalculateContentTypeFn func(repo testkube.Repository) (string, error)
}

func (f MockFetcher) Fetch(content *testkube.TestContent) (path string, err error) {
	if f.FetchFn == nil {
		log.Fatal("not implemented")
	}
	return f.FetchFn(content)
}

func (f MockFetcher) FetchString(str string) (path string, err error) {
	if f.FetchStringFn == nil {
		log.Fatal("not implemented")
	}
	return f.FetchStringFn(str)
}

func (f MockFetcher) FetchURI(str string) (path string, err error) {
	if f.FetchURIFn == nil {
		log.Fatal("not implemented")
	}
	return f.FetchURIFn(str)
}

func (f MockFetcher) FetchGitDir(repo *testkube.Repository) (path string, err error) {
	if f.FetchGitDirFn == nil {
		log.Fatal("not implemented")
	}
	return f.FetchGitDir(repo)
}

func (f MockFetcher) FetchGitFile(repo *testkube.Repository) (path string, err error) {
	if f.FetchGitFileFn == nil {
		log.Fatal("not implemented")
	}
	return f.FetchGitFileFn(repo)
}

func (f MockFetcher) FetchGit(repo *testkube.Repository) (path string, err error) {
	if f.FetchGitFn == nil {
		log.Fatal("not implemented")
	}
	return f.FetchGit(repo)
}

func (f MockFetcher) CalculateGitContentType(repo testkube.Repository) (string, error) {
	if f.FetchCalculateContentTypeFn == nil {
		log.Fatal("not implemented")
	}
	return f.FetchCalculateContentTypeFn(repo)
}
