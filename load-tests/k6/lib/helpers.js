import { check, sleep } from 'k6';
import http from 'k6/http';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.5.0/index.js';

export const INGESTION_URL = __ENV.INGESTION_URL || 'http://localhost:8081';
export const GATEWAY_URL = __ENV.GATEWAY_URL || 'http://localhost:8080';
export const API_KEY = __ENV.API_KEY || 'dev-api-key-001';

const banks = ['HDFC', 'ICICI', 'SBI', 'AXIS', 'KOTAK', 'PNB', 'BOB', 'YES'];
const vpas = ['paytm', 'ybl', 'okaxis', 'ibl'];

export function randomTxn() {
  const id = uuidv4();
  return {
    transaction_id: `TXN-${id}`,
    idempotency_key: `idem-${id}`,
    amount_paise: 100 + Math.floor(Math.random() * 500000),
    currency: 'INR',
    payer_vpa: `user${Math.floor(Math.random() * 99999)}@${vpas[Math.floor(Math.random() * vpas.length)]}`,
    payee_vpa: `merchant${Math.floor(Math.random() * 9999)}@${vpas[Math.floor(Math.random() * vpas.length)]}`,
    bank_code: banks[Math.floor(Math.random() * banks.length)],
    merchant_id: `MRC-${Math.floor(Math.random() * 99999)}`,
    txn_type: 'P2M',
    device_fingerprint: uuidv4(),
    created_at: new Date().toISOString(),
  };
}

export function postTransaction(baseUrl, useGateway = false) {
  const txn = randomTxn();
  const url = useGateway ? `${baseUrl}/v1/transactions` : `${baseUrl}/v1/transactions`;
  const headers = {
    'Content-Type': 'application/json',
    'Idempotency-Key': txn.idempotency_key,
    'X-Correlation-Id': uuidv4(),
  };
  if (useGateway) headers['Authorization'] = `Bearer ${API_KEY}`;

  const res = http.post(url, JSON.stringify(txn), { headers, tags: { name: 'ingest' } });
  check(res, {
    'status accepted': (r) => r.status === 202 || r.status === 409,
  });
  return res;
}

export function defaultOptions(vus, duration) {
  return {
    vus,
    duration,
    thresholds: {
      http_req_failed: ['rate<0.15'],
      http_req_duration: ['p(95)<2000'],
    },
  };
}
