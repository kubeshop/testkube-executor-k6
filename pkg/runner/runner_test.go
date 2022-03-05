package runner

import (
	"io/ioutil"
	"testing"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/stretchr/testify/assert"
)

func TestRunFiles(t *testing.T) {

	t.Run("Run k6 with simple script", func(t *testing.T) {
		// given
		k6_script, err := ioutil.ReadFile("k6-test-script.js")
		if err != nil {
			assert.FailNow(t, "Unable to read k6 test script.")
		}

		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = testkube.NewStringTestContent(string(k6_script))

		// when
		result, err := runner.Run(*execution)

		// then
		assert.NoError(t, err)
		assert.Equal(t, result.Status, testkube.ExecutionStatusSuccess)
	})

	t.Run("Run k6 with arguments and simple script", func(t *testing.T) {
		// given
		k6_script, err := ioutil.ReadFile("k6-test-script.js")
		if err != nil {
			assert.FailNow(t, "Unable to read k6 test script.")
		}

		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = testkube.NewStringTestContent(string(k6_script))
		execution.Args = []string{"--vus", "2", "--duration", "1s"}

		// when
		result, err := runner.Run(*execution)

		// then
		assert.NoError(t, err)
		assert.Equal(t, result.Status, testkube.ExecutionStatusSuccess)
	})

	t.Run("Run k6 with ENV variables and script", func(t *testing.T) {
		// given
		k6_script, err := ioutil.ReadFile("k6-env-test-script.js")
		if err != nil {
			assert.FailNow(t, "Unable to read k6 test script.")
		}

		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = testkube.NewStringTestContent(string(k6_script))
		execution.Envs = map[string]string{"TARGET_HOSTNAME": "kubeshop.github.io"}

		// when
		result, err := runner.Run(*execution)

		// then
		assert.NoError(t, err)
		assert.Equal(t, result.Status, testkube.ExecutionStatusSuccess)
	})
}

func TestRunDirs(t *testing.T) {
	t.Run("Run k6 from directory with script argument", func(t *testing.T) {
		// given
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = &testkube.TestContent{
			Type_: string(testkube.TestContentTypeGitDir),
			Repository: &testkube.Repository{
				Uri:    "https://github.com/kubeshop/testkube-executor-k6.git",
				Branch: "main",
				Path:   "examples",
			},
		}
		execution.Args = []string{"--duration", "1s", "k6-test-script.js"}

		// when
		result, err := runner.Run(*execution)

		// then
		assert.NoError(t, err)
		assert.Equal(t, result.Status, testkube.ExecutionStatusSuccess)
	})
}

func TestRunErrors(t *testing.T) {

	t.Run("Run k6 with no script", func(t *testing.T) {
		// given
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = testkube.NewStringTestContent("")

		// when
		result, err := runner.Run(*execution)

		// then
		assert.NoError(t, err)
		assert.Equal(t, result.Status, testkube.ExecutionStatusError)
	})

	t.Run("Run k6 with invalid arguments", func(t *testing.T) {
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = testkube.NewStringTestContent("")
		execution.Args = []string{"--vues", "2", "--duration", "5"}

		// when
		result, err := runner.Run(*execution)

		// then
		assert.NoError(t, err)
		assert.Equal(t, result.Status, testkube.ExecutionStatusError)
	})

	t.Run("Run k6 from directory with missing script arg", func(t *testing.T) {
		// given
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = &testkube.TestContent{
			Type_: string(testkube.TestContentTypeGitDir),
			Repository: &testkube.Repository{
				Uri:    "https://github.com/kubeshop/testkube-executor-k6.git",
				Branch: "main",
				Path:   "examples",
			},
		}
		execution.Args = []string{}

		// when
		result, err := runner.Run(*execution)

		// then
		assert.NoError(t, err)
		assert.Equal(t, result.Status, testkube.ExecutionStatusError)
	})
}
