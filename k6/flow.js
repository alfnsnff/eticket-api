import http from 'k6/http';
import { check, group, sleep } from 'k6';

export { options } from './scenarios/load.js';

const BASE_URL = 'http://206.189.157.95:8080/api';
const AVAILABLE_SCHEDULES = [1, 2, 3, 4];
const CLASS_ID = 1;

export default function () {
  const vuId = __VU;
  const iteration = __ITER;

  group('Get schedules', () => {
    const res = http.get(`${BASE_URL}/v1/schedules`, { tags: { name: '/schedules' } });
    const ok = check(res, { 'load schedules success': (r) => r.status === 200 });
    if (!ok) console.error(JSON.stringify({ vu: vuId, iter: iteration, endpoint: 'GET /v1/schedules', status: res.status, body: res.body }));
    sleep(1);
  });

  let selectedSchedule;
  group('Get Schedule By ID', () => {
    selectedSchedule = AVAILABLE_SCHEDULES[Math.floor(Math.random() * AVAILABLE_SCHEDULES.length)];
    const res = http.get(`${BASE_URL}/v1/schedule/${selectedSchedule}`, { tags: { name: '/schedule/:id' } });
    const ok = check(res, { 'load class success': (r) => r.status === 200 });
    if (!ok) console.error(JSON.stringify({ vu: vuId, iter: iteration, endpoint: `GET /v1/schedule/${selectedSchedule}`, status: res.status, body: res.body }));
    sleep(1);
  });

  let sessionId;
  group('Claim Lock', () => {
    const payload = JSON.stringify({ schedule_id: selectedSchedule, items: [{ class_id: CLASS_ID, quantity: 1 }] });
    const res = http.post(`${BASE_URL}/v1/claim/lock`, payload, {
      headers: { 'Content-Type': 'application/json' },
      tags: { name: '/claim/lock' },
    });
    const ok = check(res, { 'claim lock success': (r) => r.status === 201 });
    if (!ok) {
      console.error(JSON.stringify({ vu: vuId, iter: iteration, endpoint: 'POST /v1/claim/lock', status: res.status, body: res.body }));
      return;
    }
    sessionId = JSON.parse(res.body).data.session_id;
    sleep(0.5);
  });

  let orderId;
  group('Claim Entry', () => {
    const payload = JSON.stringify({
      customer_name: `LoadTest${vuId}${iteration}`,
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

    const res = http.post(`${BASE_URL}/v1/claim/entry/${sessionId}`, payload, {
      headers: { 'Content-Type': 'application/json' },
      tags: { name: '/claim/entry' },
    });
    const ok = check(res, { 'claim entry success': (r) => r.status === 200 });
    if (!ok) {
      console.error(JSON.stringify({ vu: vuId, iter: iteration, endpoint: `POST /v1/claim/entry/${sessionId}`, status: res.status, body: res.body }));
      return;
    }
    orderId = JSON.parse(res.body).data.order_id;
    sleep(1);
  });

  group('Payment Callback', () => {
    const payload = JSON.stringify({
      reference: `LT_${Date.now()}_${vuId}_${iteration}`,
      merchant_ref: orderId,
      status: 'PAID',
      amount: 100000,
      payment_method: 'BRIVA',
      signature: `SIGN_${Date.now()}_${vuId}`,
    });

    const res = http.post(`${BASE_URL}/v1/payment/callback`, payload, {
      headers: { 'Content-Type': 'application/json' },
      tags: { name: '/payment/callback' },
    });
    const ok = check(res, { 'payment success': (r) => r.status === 200 });
    if (!ok) console.error(JSON.stringify({ vu: vuId, iter: iteration, endpoint: 'POST /v1/payment/callback', status: res.status, body: res.body }));
    sleep(1);
  });

  group('Get Booking', () => {
    const res = http.get(`${BASE_URL}/v1/booking/order/${orderId}`, { tags: { name: 'get_booking' } });
    const ok = check(res, { 'get booking success': (r) => r.status === 200 });
    if (!ok) console.error(JSON.stringify({ vu: vuId, iter: iteration, endpoint: `GET /v1/booking/order/${orderId}`, status: res.status, body: res.body }));
    sleep(Math.random() * 3 + 2);
  });
}