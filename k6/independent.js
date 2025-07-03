import http from 'k6/http';
import { check, group } from 'k6';

export { options } from './scenarios/load.js';

const BASE_URL = 'http://206.189.157.95:8080/api';
const AVAILABLE_SCHEDULES = [1, 2, 3, 4];
const CLASS_ID = 1;
const SESSION_ID = '5064e05a-b9aa-4fa7-b70b-35a4aa6c423f'; // sesuaikan jika perlu
const ORDER_ID = 'PH2-20250703052620882469231-a57b61af95515f93'; // sesuaikan jika perlu

export default function () {
  const vuId = __VU;
  const iteration = __ITER;

  group('Get Schedules', () => {
    const res = http.get(`${BASE_URL}/v1/schedules`, {
      tags: { name: '/schedules' },
    });
    const ok = check(res, { 'load schedules success': (r) => r.status === 200 });
    if (!ok) console.error(JSON.stringify({ vu: vuId, iter: iteration, endpoint: 'GET /v1/schedules', status: res.status, body: res.body }));
  });

  group('Get Schedule by ID', () => {
    const sid = AVAILABLE_SCHEDULES[Math.floor(Math.random() * AVAILABLE_SCHEDULES.length)];
    const res = http.get(`${BASE_URL}/v1/schedule/${sid}`, {
      tags: { name: '/schedule/:id' },
    });
    const ok = check(res, { 'load class success': (r) => r.status === 200 });
    if (!ok) console.error(JSON.stringify({ vu: vuId, iter: iteration, endpoint: `GET /v1/schedule/${sid}`, status: res.status, body: res.body }));
  });

  group('Claim Lock', () => {
    const sid = AVAILABLE_SCHEDULES[Math.floor(Math.random() * AVAILABLE_SCHEDULES.length)];
    const payload = JSON.stringify({
      schedule_id: sid,
      items: [{ class_id: CLASS_ID, quantity: 1 }],
    });
    const res = http.post(`${BASE_URL}/v1/claim/lock`, payload, {
      headers: { 'Content-Type': 'application/json' },
      tags: { name: '/claim/lock' },
    });
    const ok = check(res, { 'claim lock success': (r) => r.status === 201 });
    if (!ok) console.error(JSON.stringify({ vu: vuId, iter: iteration, endpoint: 'POST /v1/claim/lock', status: res.status, body: res.body }));
  });

  group('Claim Entry', () => {
    const payload = JSON.stringify({
      customer_name: `TestUser`,
      id_type: 'ktp',
      id_number: `${Math.floor(Math.random() * 9e15) + 1e15}`,
      phone_number: `08${String(Math.floor(Math.random() * 1e9)).padStart(9, '0')}`,
      email: `lt${vuId}${iteration}@test.com`,
      payment_method: 'BRIVA',
      ticket_data: [{
        class_id: CLASS_ID,
        passenger_name: `Pass${vuId}${iteration}`,
        passenger_age: Math.floor(Math.random() * 50) + 18,
        passenger_gender: Math.random() > 0.5 ? 'male' : 'female',
        id_type: 'ktp',
        id_number: `${Math.floor(Math.random() * 9e15) + 1e15}`,
        address: `Addr ${vuId} ${iteration}`,
      }],
    });
    const res = http.post(`${BASE_URL}/v1/claim/entry/${SESSION_ID}`, payload, {
      headers: { 'Content-Type': 'application/json' },
      tags: { name: '/claim/entry/:sessionId' },
    });
    const ok = check(res, { 'claim entry success': (r) => r.status === 200 });
    if (!ok) console.error(JSON.stringify({ vu: vuId, iter: iteration, endpoint: `POST /v1/claim/entry/${SESSION_ID}`, status: res.status, body: res.body }));
  });

  group('Get Booking by Order ID', () => {
    const res = http.get(`${BASE_URL}/v1/booking/order/${ORDER_ID}`, {
      tags: { name: '/booking/order/:id' },
    });
    const ok = check(res, { 'get booking success': (r) => r.status === 200 });
    if (!ok) console.error(JSON.stringify({ vu: vuId, iter: iteration, endpoint: `GET /v1/booking/order/${ORDER_ID}`, status: res.status, body: res.body }));
  });

  group('Payment Callback', () => {
    const payload = JSON.stringify({
      reference: `LT_${Date.now()}_${vuId}_${iteration}`,
      merchant_ref: ORDER_ID,
      status: 'PAID',
      amount: 100000,
      payment_method: 'BRIVA',
      signature: `SIGN_${Date.now()}_${vuId}`,
    });
    const res = http.post(`${BASE_URL}/v1/payment/callback`, payload, {
      headers: { 'Content-Type': 'application/json' },
      tags: { name: '/payment/callback' },
    });
    const ok = check(res, { 'payment callback success': (r) => r.status === 200 });
    if (!ok) console.error(JSON.stringify({ vu: vuId, iter: iteration, endpoint: 'POST /v1/payment/callback', status: res.status, body: res.body }));
  });
}
