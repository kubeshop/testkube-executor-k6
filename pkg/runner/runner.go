package runner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor"
	"github.com/kubeshop/testkube/pkg/executor/output"
)

type Params struct {
	Datadir string // RUNNER_DATADIR
}

func NewRunner() *K6Runner {
	params := Params{
		Datadir: os.Getenv("RUNNER_DATADIR"),
	}

	runner := &K6Runner{
		Params: params,
	}

	return runner
}

type K6Runner struct {
	Params Params
}

func (r *K6Runner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {
	// check that the datadir exists
	_, err = os.Stat(r.Params.Datadir)
	if errors.Is(err, os.ErrNotExist) {
		return result, err
	}

	args := []string{"run"}

	// convert executor env variables to k6 env variables
	for key, value := range execution.Envs {
		env_var := fmt.Sprintf("%s=%s", key, value)
		args = append(args, "-e", env_var)
	}

	// pass additional executor arguments/flags to k6
	args = append(args, execution.Args...)

	var directory string

	// in case of a test file execution we will pass the
	// file path as final parameter to k6
	if execution.Content.IsFile() {
		args = append(args, "test-content")
		directory = r.Params.Datadir
	}

	// in case of Git directory we will run k6 here and
	// use the last argument as test file
	if execution.Content.IsDir() {
		directory = filepath.Join(r.Params.Datadir, "repo")

		// sanity checking for test script
		script_file := filepath.Join(directory, args[len(args)-1])
		file_info, err := os.Stat(script_file)
		if errors.Is(err, os.ErrNotExist) || file_info.IsDir() {
			return result.Err(fmt.Errorf("k6 test script %s not found", script_file)), nil
		}
	}

	output.PrintEvent("Running", directory, "k6", args)
	output, err := executor.Run(directory, "k6", args...)
	if err != nil {
		return result.Err(err), nil
	}

	return testkube.ExecutionResult{
		Status:     testkube.StatusPtr(testkube.SUCCESS_ExecutionStatus),
		Output:     string(output),
		OutputType: "text/plain",
		Steps: []testkube.ExecutionStepResult{
			{
				// use the scenario name and description here
				Name:     parseScenarioName(string(output)),
				Duration: parseScenarioDuration(string(output)),
				Status:   parseScenarioStatus(string(output)),
			},
		},
	}, nil
}

func parseScenarioName(summary string) string {
	lines := splitSummaryBody(summary)
	var name string

	for i, line := range lines {
		if strings.Contains(line, "scenario") {
			// take next line and trim leading spaces
			name = strings.TrimSpace(lines[i+1])
			break
		}
	}

	return name
}

func parseScenarioDuration(summary string) string {
	lines := splitSummaryBody(summary)
	var duration string

	for i, line := range lines {
		if strings.Contains(line, "running") {
			// take next line and trim leading spaces
			metrics := strings.Split(lines[i+1], " ")
			duration = metrics[6]
			break
		}
	}

	return duration
}

func parseScenarioStatus(summary string) string {
	return string(testkube.SUCCESS_ExecutionStatus)
}

func splitSummaryBody(summary string) []string {
	return strings.Split(summary, "\n")
}
