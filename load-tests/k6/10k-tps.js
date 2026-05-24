import { sleep } from 'k6';
import { postTransaction, INGESTION_URL } from './lib/helpers.js';

// Target ~10K TPS — tune VUs based on machine (start 200-500 VUs)
export const options = {
  scenarios: {
    sustained: {
      executor: 'constant-arrival-rate',
      rate: 10000,
      timeUnit: '1s',
      duration: '3m',
      preAllocatedVUs: 500,
      maxVUs: 2000,
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.20'],
    http_req_duration: ['p(95)<3000'],
  },
};

export default function () {
  postTransaction(INGESTION_URL);
  sleep(0.001);
}
