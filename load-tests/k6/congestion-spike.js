import { sleep } from 'k6';
import { postTransaction, GATEWAY_URL, defaultOptions } from './lib/helpers.js';

export const options = {
  stages: [
    { duration: '30s', target: 100 },
    { duration: '1m', target: 600 },
    { duration: '30s', target: 100 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.2'],
    http_req_duration: ['p(95)<2500'],
  },
};

export default function () {
  postTransaction(GATEWAY_URL, true);
  sleep(0.03);
}
