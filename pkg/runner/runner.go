package runner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

	args := []string{"run"}

	// convert executor env variables to k6 env variables
	for key, value := range execution.Envs {
		env_var := fmt.Sprintf("%s=%s", key, value)
		args = append(args, "-e", env_var)
	}

	// pass additional executor arguments/flags to k6
	args = append(args, execution.Args...)

	// in case of a test file execution we will pass the
	// file path as final parameter to k6
	if execution.Content.IsFile() {
		args = append(args, path)
	}

	// in case of Git directory we will run k6 here
	directory := ""
	if execution.Content.IsDir() {
		directory = path

		// sanity checking
		// the last argument needs to be an existing file
		script_file := filepath.Join(directory, args[len(args)-1])
		file_info, err := os.Stat(script_file)
		if errors.Is(err, os.ErrNotExist) || file_info.IsDir() {
			return result, fmt.Errorf("k6 script %s not found", script_file)
		}
	}

	output.PrintEvent("Running k6", args)
	output, err := executor.Run(directory, "k6", args...)
	if err != nil {
		return result.Err(err), nil
	}

	return testkube.ExecutionResult{
		Status: testkube.StatusPtr(testkube.SUCCESS_ExecutionStatus),
		Output: string(output),
	}, nil
}
