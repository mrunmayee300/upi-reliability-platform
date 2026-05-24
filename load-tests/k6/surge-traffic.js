import { sleep } from 'k6';
import { postTransaction, INGESTION_URL } from './lib/helpers.js';

export const options = {
  stages: [
    { duration: '1m', target: 50 },
    { duration: '30s', target: 800 },
    { duration: '2m', target: 800 },
    { duration: '30s', target: 50 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.25'],
  },
};

export default function () {
  postTransaction(INGESTION_URL);
  sleep(0.01);
}
