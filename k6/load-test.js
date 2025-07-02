import { check, sleep } from 'k6';
import http from 'k6/http';

export { options } from './scenarios/spike.js';

const BASE_URL = 'http://localhost:8080/api';
const AVAILABLE_SCHEDULES = [9, 10, 11, 12];
const CLASS_ID = 4;

export default function () {
  const vuId = __VU;
  const iteration = __ITER;

  // 1. Load schedules
  const scheduleResponse = http.get(`${BASE_URL}/v1/schedules`);
  const scheduleSuccess = check(scheduleResponse, {
    'load schedules success': (r) => r.status === 200,
  });
  if (!scheduleSuccess) {
    console.error(JSON.stringify({
      vu: vuId,
      iter: iteration,
      endpoint: 'GET /v1/schedules',
      status: scheduleResponse.status,
      body: scheduleResponse.body,
    }));
  }
  sleep(1);

  // 2. Load class
  const selectedSchedule = AVAILABLE_SCHEDULES[Math.floor(Math.random() * AVAILABLE_SCHEDULES.length)];
  const classResponse = http.get(`${BASE_URL}/v1/schedule/${selectedSchedule}`);
  const classSuccess = check(classResponse, {
    'load class success': (r) => r.status === 200,
  });
  if (!classSuccess) {
    console.error(JSON.stringify({
      vu: vuId,
      iter: iteration,
      endpoint: `GET /v1/schedule/${selectedSchedule}`,
      status: classResponse.status,
      body: classResponse.body,
    }));
  }
  sleep(1);

  // 3. Claim Lock
  const lockData = {
    schedule_id: selectedSchedule,
    items: [{ class_id: CLASS_ID, quantity: 1 }],
  };
  const lockResponse = http.post(`${BASE_URL}/v1/claim/lock`, JSON.stringify(lockData), {
    headers: { 'Content-Type': 'application/json' },
  });
  const lockSuccess = check(lockResponse, {
    'claim lock success': (r) => r.status === 201,
  });
  if (!lockSuccess) {
    console.error(JSON.stringify({
      vu: vuId,
      iter: iteration,
      endpoint: 'POST /v1/claim/lock',
      status: lockResponse.status,
      body: lockResponse.body,
    }));
    return;
  }

  const sessionId = JSON.parse(lockResponse.body).data.session_id;
  sleep(0.5);

  // 4. Claim Entry
  const entryData = {
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
  };

  const entryResponse = http.post(`${BASE_URL}/v1/claim/entry/${sessionId}`, JSON.stringify(entryData), {
    headers: { 'Content-Type': 'application/json' },
  });

  const entrySuccess = check(entryResponse, {
    'claim entry success': (r) => r.status === 200,
  });
  if (!entrySuccess) {
    console.error(JSON.stringify({
      vu: vuId,
      iter: iteration,
      endpoint: `POST /v1/claim/entry/${sessionId}`,
      status: entryResponse.status,
      body: entryResponse.body,
    }));
    return;
  }

  const orderId = JSON.parse(entryResponse.body).data.order_id;
  sleep(1);

  // 5. Payment Callback
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
  const paymentSuccess = check(paymentResponse, {
    'payment success': (r) => r.status === 200,
  });
  if (!paymentSuccess) {
    console.error(JSON.stringify({
      vu: vuId,
      iter: iteration,
      endpoint: 'POST /v1/payment/callback',
      status: paymentResponse.status,
      body: paymentResponse.body,
    }));
  }
  sleep(1);

  // 6. Get Booking
  const bookingResponse = http.get(`${BASE_URL}/v1/booking/order/${orderId}`);
  const bookingSuccess = check(bookingResponse, {
    'get booking success': (r) => r.status === 200,
  });
  if (!bookingSuccess) {
    console.error(JSON.stringify({
      vu: vuId,
      iter: iteration,
      endpoint: `GET /v1/booking/order/${orderId}`,
      status: bookingResponse.status,
      body: bookingResponse.body,
    }));
  }

  sleep(Math.random() * 3 + 2);
}
