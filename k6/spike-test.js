import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 0 },   // Awal idle
    { duration: '10s', target: 50 },  // Spike mendadak
    { duration: '30s', target: 50 },  // Tahan di spike
    { duration: '10s', target: 0 },   // Turunkan cepat
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'],
    http_req_failed: ['rate<0.05'],
  },
};

const BASE_URL = 'http://localhost:8080/api';
const AVAILABLE_SCHEDULES = [9, 10, 11, 12, 13];
const SCHEDULE_NAMES = {
  9: 'Jakarta-Bandung',
  10: 'Jakarta-Surabaya',
  11: 'Bandung-Yogyakarta',
  12: 'Surabaya-Malang',
  13: 'Jakarta-Semarang',
};

export default function () {
  const vuId = __VU;
  const iteration = __ITER;

  const scheduleResponse = http.get(`${BASE_URL}/v1/schedules`);
  check(scheduleResponse, { 'load schedules': (r) => r.status === 200 });
  sleep(1);

  const classResponse = http.get(`${BASE_URL}/v1/classes`);
  check(classResponse, { 'load classes': (r) => r.status === 200 });
  sleep(1);

  const selectedSchedule = AVAILABLE_SCHEDULES[Math.floor(Math.random() * AVAILABLE_SCHEDULES.length)];
  const lockData = {
    schedule_id: selectedSchedule,
    items: [{ class_id: 4, quantity: 1 }],
  };

  const lockResponse = http.post(`${BASE_URL}/v1/claim/lock`, JSON.stringify(lockData), {
    headers: { 'Content-Type': 'application/json' },
  });

  if (lockResponse.status !== 201) return;

  let sessionId;
  try {
    sessionId = JSON.parse(lockResponse.body).data.session_id;
  } catch (_) {
    return;
  }

  sleep(0.5);

  const entryData = {
    customer_name: `LoadTest${vuId}${iteration}`,
    id_type: 'ktp',
    id_number: `${Math.floor(Math.random() * 9000000000000000) + 1000000000000000}`,
    phone_number: `08${String(Math.floor(Math.random() * 1000000000)).padStart(9, '0')}`,
    email: `lt${vuId}${iteration}@test.com`,
    payment_method: 'BRIVA',
    ticket_data: [{
      class_id: 4,
      passenger_name: `Pass${vuId}${iteration}`,
      passenger_age: Math.floor(Math.random() * 50) + 18,
      passenger_gender: Math.random() > 0.5 ? 'male' : 'female',
      id_type: 'ktp',
      id_number: `${Math.floor(Math.random() * 9000000000000000) + 1000000000000000}`,
      address: `Addr ${vuId} ${iteration}`,
    }]
  };

  const entryResponse = http.post(`${BASE_URL}/v1/claim/entry/${sessionId}`, JSON.stringify(entryData), {
    headers: { 'Content-Type': 'application/json' },
  });

  if (entryResponse.status !== 200) return;

  let orderId;
  try {
    orderId = JSON.parse(entryResponse.body).data.order_id;
  } catch (_) {
    return;
  }

  sleep(1);

  const paymentData = {
    reference: `LT_${Date.now()}_${vuId}_${iteration}`,
    merchant_ref: orderId,
    status: 'PAID',
    amount: 100000,
    payment_method: 'BRIVA',
    signature: `SIGN_${Date.now()}_${vuId}`,
  };

  const paymentResponse = http.post(`${BASE_URL}/v1/payment/callback`, JSON.stringify(paymentData), {
    headers: { 'Content-Type': 'application/json' },
  });

  check(paymentResponse, { 'payment success': (r) => r.status === 200 });

  sleep(Math.random() * 2 + 1); // think time 1–3s
}
