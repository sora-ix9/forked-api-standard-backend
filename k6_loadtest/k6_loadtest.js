import http from 'k6/http';
import { check, sleep } from 'k6';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";
import { textSummary } from "https://jslib.k6.io/k6-summary/0.0.1/index.js";

export const options = {
  stages: [
    { duration: '30s', target: 20 }, // Ramp up to 20 users
    { duration: '1m', target: 20 },  // Stay at 20 users
    { duration: '10s', target: 0 },  // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests must complete below 500ms
  },
};

const BASE_URL = 'http://localhost:1323';

export default function () {
  // 1. Health Check
  const res = http.get(`${BASE_URL}/`);
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response body contains online': (r) => r.body.includes('Online'),
  });

  // 2. Login (Example - Uncomment and adjust with valid data to test)
  /*
  const loginPayload = JSON.stringify({
    email: 'test@example.com',
    password: 'password123',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const loginRes = http.post(`${BASE_URL}/user/login`, loginPayload, params);
  
  check(loginRes, {
    'login status is 200': (r) => r.status === 200,
    'has token': (r) => r.json('token') !== undefined,
  });
  */

  sleep(1);
}

export function handleSummary(data) {
  return {
    "k6_loadtest/summary.html": htmlReport(data),
    stdout: textSummary(data, { indent: " ", enableColors: true }),
  };
}
