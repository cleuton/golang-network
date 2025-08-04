import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 50 },    // ramp up to 50 users
    { duration: '1m', target: 50 },     // stay at 50 users
    { duration: '30s', target: 0 },     // ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],   // 95% of requests should be below 500ms
  },
};

export default function () {
  const res = http.get('http://localhost:8888/nextid');
  check(res, {
    'status is 200': (r) => r.status === 200,
    'has id': (r) => r.json().id !== undefined,
  });
  
  // Thinking time: 50ms entre requisições
  sleep(0.05);
}