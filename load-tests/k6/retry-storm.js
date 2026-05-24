import { sleep } from 'k6';
import { postTransaction, INGESTION_URL, defaultOptions } from './lib/helpers.js';

export const options = {
  ...defaultOptions(120, '2m'),
  thresholds: {
    http_req_failed: ['rate<0.25'],
    http_req_duration: ['p(95)<3000'],
  },
};

export default function () {
  postTransaction(INGESTION_URL);
  sleep(0.01);
}
