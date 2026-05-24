import { sleep } from 'k6';
import { postTransaction, GATEWAY_URL, defaultOptions } from './lib/helpers.js';

export const options = defaultOptions(100, '3m');

export default function () {
  postTransaction(GATEWAY_URL, true);
  sleep(0.02);
}
