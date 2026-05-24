import { sleep } from 'k6';
import { postTransaction, INGESTION_URL, defaultOptions } from './lib/helpers.js';

export const options = defaultOptions(80, '3m');

export default function () {
  postTransaction(INGESTION_URL);
  sleep(0.02);
}
