![Testkube Logo](https://raw.githubusercontent.com/kubeshop/testkube/main/assets/testkube-color-gray.png)

# Welcome to TestKube k6 Executor

TestKube k6 Executor is a test executor to run k6 load tests with [TestKube](https://testkube.io).  

## Usage

You need to register and deploy the executor in your cluster. Additionally, you may deploy InfluxDB as well as Grafana if you need detailed performance data from your tests.
```bash
kubectl apply -f examples/k6-executor.yaml

# see https://k6.io/docs/results-visualization/influxdb-+-grafana/
kubectl apply -f examples/k6-influxdb-grafana.yaml
```

Have a look at the [k6 documentation](https://k6.io/docs/getting-started/running-k6/) for details on writing tests. Here is a simple example script:
```javascript
import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  http.get('https://kubeshop.github.io/testkube/');
  sleep(1);
}
```

Issue the following commands to create and start the script:
```bash
kubectl testkube tests create --file examples/k6-test-script.js --type "k6/script" --name k6-test-script
kubectl testkube tests run k6-test-script --watch k6-test-script
```

## Examples

TODO: add more examples

# Issues and enchancements 

Please follow the main [TestKube repository](https://github.com/kubeshop/testkube) for reporting any [issues](https://github.com/kubeshop/testkube/issues) or [discussions](https://github.com/kubeshop/testkube/discussions)

# Testkube 

For more info go to [main testkube repo](https://github.com/kubeshop/testkube)

![Release](https://img.shields.io/github/v/release/kubeshop/testkube) [![Releases](https://img.shields.io/github/downloads/kubeshop/testkube/total.svg)](https://github.com/kubeshop/testkube/tags?label=Downloads) ![Go version](https://img.shields.io/github/go-mod/go-version/kubeshop/testkube)

![Docker builds](https://img.shields.io/docker/automated/kubeshop/testkube-api-server) ![Code build](https://img.shields.io/github/workflow/status/kubeshop/testkube/Code%20build%20and%20checks) ![Release date](https://img.shields.io/github/release-date/kubeshop/testkube)

![Twitter](https://img.shields.io/twitter/follow/thekubeshop?style=social) ![Discord](https://img.shields.io/discord/884464549347074049)
 #### [Documentation](https://kubeshop.github.io/testkube) | [Discord](https://discord.gg/hfq44wtR6Q) 