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

const K6_CLOUD = "cloud"
const K6_RUN = "run"
const K6_SCRIPT = "script"

func (r *K6Runner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {
	// check that the datadir exists
	_, err = os.Stat(r.Params.Datadir)
	if errors.Is(err, os.ErrNotExist) {
		return result, err
	}

	args := []string{}

	k6TestType := strings.Split(execution.TestType, "/")
	if len(k6TestType) != 2 {
		return result.Err(fmt.Errorf("invalid test type %s", execution.TestType)), nil
	}

	k6Subtype := k6TestType[1]
	if k6Subtype == K6_SCRIPT || k6Subtype == K6_RUN {
		args = append(args, K6_RUN)
	} else if k6Subtype == K6_CLOUD {
		args = append(args, K6_CLOUD)
	} else {
		return result.Err(fmt.Errorf("unsupported test type %s", execution.TestType)), nil
	}

	// convert executor env variables to k6 env variables
	for key, value := range execution.Envs {
		if key == "K6_CLOUD_TOKEN" {
			// set as OS environment variable
			os.Setenv(key, value)
		} else {
			// pass to k6 using -e option
			env := fmt.Sprintf("%s=%s", key, value)
			args = append(args, "-e", env)
		}
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
		scriptFile := filepath.Join(directory, args[len(args)-1])
		fileInfo, err := os.Stat(scriptFile)
		if errors.Is(err, os.ErrNotExist) || fileInfo.IsDir() {
			return result.Err(fmt.Errorf("k6 test script %s not found", scriptFile)), nil
		}
	}

	output.PrintEvent("Running", directory, "k6", args)
	output, err := executor.Run(directory, "k6", args...)
	return finalExecutionResult(string(output), err), nil
}

// finalExecutionResult processes the output of the test run
func finalExecutionResult(output string, err error) (result testkube.ExecutionResult) {
	switch {
	case err == nil && areChecksSuccessful(output):
		result.Status = testkube.ExecutionStatusPassed
	case err == nil && !areChecksSuccessful(output):
		result.Status = testkube.ExecutionStatusFailed
		result.ErrorMessage = "some checks have failed"
	case err != nil && strings.Contains(err.Error(), "exit status 99"):
		// tests have run, but some checks + thresholds have failed
		result.Status = testkube.ExecutionStatusFailed
		result.ErrorMessage = "some thresholds have failed"
	default:
		// k6 was unable to run at all
		result.Status = testkube.ExecutionStatusFailed
		result.ErrorMessage = err.Error()
		return result
	}

	// always set these, no matter if error or success
	result.Output = output
	result.OutputType = "text/plain"

	result.Steps = []testkube.ExecutionStepResult{}
	for _, name := range parseScenarioNames(output) {
		result.Steps = append(result.Steps, testkube.ExecutionStepResult{
			// use the scenario name with description here
			Name:     name,
			Duration: parseScenarioDuration(output, splitScenarioName(name)),

			// currently there is no way to extract individual scenario status
			Status: string(testkube.PASSED_ExecutionStatus),
		})
	}

	return result
}

// areChecksSuccessful verifies the summary at the end of the execution to see if any of the checks failed
func areChecksSuccessful(summary string) bool {
	lines := splitSummaryBody(summary)

	for _, line := range lines {
		if !strings.Contains(line, "checks") {
			continue
		}
		if strings.Contains(line, "100.00%") {
			return true
		}
		return false
	}

	return true
}

func parseScenarioNames(summary string) []string {
	lines := splitSummaryBody(summary)
	names := []string{}

	for _, line := range lines {
		if strings.Contains(line, "* ") {
			name := strings.TrimLeft(strings.TrimSpace(line), "* ")
			names = append(names, name)
		}
	}

	return names
}

func parseScenarioDuration(summary string, name string) string {
	lines := splitSummaryBody(summary)

	var duration string
	for _, line := range lines {
		if strings.Contains(line, name) && strings.Contains(line, "[ 100% ]") {
			index := strings.Index(line, "]") + 1
			line = strings.TrimSpace(line[index:])
			line = strings.ReplaceAll(line, "  ", " ")

			// take next line and trim leading spaces
			metrics := strings.Split(line, " ")
			duration = metrics[2]
			break
		}
	}

	return duration
}

func splitScenarioName(name string) string {
	return strings.Split(name, ":")[0]
}

func splitSummaryBody(summary string) []string {
	return strings.Split(summary, "\n")
}
