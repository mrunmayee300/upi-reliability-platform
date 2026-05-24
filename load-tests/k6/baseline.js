import { sleep } from 'k6';
import { postTransaction, INGESTION_URL, defaultOptions } from './lib/helpers.js';

export const options = defaultOptions(50, '2m');

export default function () {
  postTransaction(INGESTION_URL);
  sleep(0.05);
}
