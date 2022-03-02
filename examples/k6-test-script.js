import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
  insecureSkipTLSVerify: true,
  thresholds: {
      'http_req_duration{kind:html}': ['avg<=250', 'p(95)<500'],
  }
};

export default function () {
  check(http.get('https://kubeshop.github.io/testkube/', {
      tags: {'kind': 'html'},
  }), {
      "status is 200": (res) => res.status === 200,
  });
  sleep(1);
}
