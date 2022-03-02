package runner

import (
	"fmt"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor"
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

	output.PrintEvent("Created content path", path)

	if !execution.Content.IsFile() {
		return result, testkube.ErrTestContentTypeNotFile
	}

	args := []string{"run"}

	// convert executor env variables to k6 env variables
	for key, value := range execution.Envs {
		env_var := fmt.Sprintf("%s=%s", key, value)
		args = append(args, "-e", env_var)
	}

	// pass additional arguments/flags to k6
	args = append(args, execution.Args...)
	args = append(args, path)

	output.PrintEvent("Running k6", args)
	output, err := executor.Run("", "k6", args...)
	if err != nil {
		return result.Err(err), nil
	}

	return testkube.ExecutionResult{
		Status: testkube.StatusPtr(testkube.SUCCESS_ExecutionStatus),
		Output: string(output),
	}, nil
}
