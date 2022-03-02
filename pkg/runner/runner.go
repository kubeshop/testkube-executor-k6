package runner

import (
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor/content"
	"github.com/kubeshop/testkube/pkg/executor/output"
)

func NewRunner() *K6Runner {
	return &K6Runner{
		Fetcher: content.NewFetcher(""),
	}
}

type K6Runner struct {
	Fetcher content.ContentFetcher
}

func (r *K6Runner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {
	path, err := r.Fetcher.Fetch(execution.Content)
	if err != nil {
		return result, err
	}

	output.PrintEvent("created content path", path)

	if execution.Content.IsFile() {
		output.PrintEvent("using file", execution)
		// TODO implement file based test content for string, git-file, file-uri
		//      or remove if not used
	}

	if execution.Content.IsDir() {
		output.PrintEvent("using dir", execution)
		// TODO implement file based test content for git-dir
		//      or remove if not used
	}

	// TODO run executor here

	// error result should be returned if something is not ok
	// return result.Err(fmt.Errorf("some test execution related error occured"))

	// TODO return ExecutionResult
	return testkube.ExecutionResult{
		Status: testkube.ExecutionStatusSuccess,
		Output: "exmaple test output",
	}, nil
}
