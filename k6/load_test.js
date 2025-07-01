import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '30s', target: 5 },   // Warm up
    { duration: '1m', target: 10 },   // Ramp up
    { duration: '2m', target: 20 },   // Sustained load
    { duration: '1m', target: 30 },   // Peak load
    { duration: '30s', target: 0 },   // Cool down
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% of requests should be below 2s
    http_req_failed: ['rate<0.1'],     // Error rate should be below 10%
    errors: ['rate<0.1'],
  },
};

const BASE_URL = 'http://localhost:8080/api'; // Ganti dengan URL API Anda

export default function () {
  // 1. Login atau Register User
  const loginResponse = login();
  if (!loginResponse) return;
  
  const authToken = loginResponse.token;
  
  // 2. Browse Events
  const events = getEvents(authToken);
  if (!events || events.length === 0) return;
  
  // 3. Create Booking
  const booking = createBooking(authToken, events[0].id);
  if (!booking) return;
  
  // 4. Process Payment
  const payment = processPayment(authToken, booking.id);
  if (!payment) return;
  
  // 5. Simulate Payment Callback (mock Tripay)
  simulatePaymentCallback(payment.reference);
  
  sleep(1); // Think time between iterations
}

function login() {
  const userData = {
    email: `testuser${Math.floor(Math.random() * 10000)}@example.com`,
    password: 'password123',
    name: `Test User ${Math.floor(Math.random() * 1000)}`,
    phone: `08${Math.floor(Math.random() * 1000000000)}`
  };
  
  // Try login first, if fails then register
  let response = http.post(`${BASE_URL}/v1/auth/login`, JSON.stringify({
    email: userData.email,
    password: userData.password
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  
  if (response.status === 401 || response.status === 404) {
    // User doesn't exist, register first
    response = http.post(`${BASE_URL}/v1/auth/register`, JSON.stringify(userData), {
      headers: { 'Content-Type': 'application/json' },
    });
    
    check(response, {
      'register successful': (r) => r.status === 201,
    }) || errorRate.add(1);
    
    if (response.status !== 201) return null;
    
    // Now login
    response = http.post(`${BASE_URL}/v1/auth/login`, JSON.stringify({
      email: userData.email,
      password: userData.password
    }), {
      headers: { 'Content-Type': 'application/json' },
    });
  }
  
  const loginSuccess = check(response, {
    'login successful': (r) => r.status === 200,
    'has auth token': (r) => r.json('token') !== undefined,
  });
  
  if (!loginSuccess) {
    errorRate.add(1);
    return null;
  }
  
  return response.json();
}

function getEvents(token) {
  const response = http.get(`${BASE_URL}/v1/events`, {
    headers: { 
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
  });
  
  const success = check(response, {
    'events retrieved': (r) => r.status === 200,
    'has events data': (r) => r.json('data') && r.json('data').length > 0,
  });
  
  if (!success) {
    errorRate.add(1);
    return null;
  }
  
  return response.json('data');
}

function createBooking(token, eventId) {
  const bookingData = {
    eventId: eventId,
    quantity: Math.floor(Math.random() * 3) + 1, // 1-3 tickets
    ticketType: 'regular',
    customerInfo: {
      name: `Customer ${Math.floor(Math.random() * 1000)}`,
      email: `customer${Math.floor(Math.random() * 10000)}@example.com`,
      phone: `08${Math.floor(Math.random() * 1000000000)}`
    }
  };
  
  const response = http.post(`${BASE_URL}/v1/bookings`, JSON.stringify(bookingData), {
    headers: { 
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
  });
  
  const success = check(response, {
    'booking created': (r) => r.status === 201,
    'has booking id': (r) => r.json('data.id') !== undefined,
  });
  
  if (!success) {
    errorRate.add(1);
    return null;
  }
  
  return response.json('data');
}

function processPayment(token, bookingId) {
  const paymentData = {
    bookingId: bookingId,
    paymentMethod: 'BRIVA', // Tripay payment method
    amount: 100000 + (Math.floor(Math.random() * 5) * 50000) // Random amount
  };
  
  const response = http.post(`${BASE_URL}/v1/payments/process`, JSON.stringify(paymentData), {
    headers: { 
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
  });
  
  const success = check(response, {
    'payment processed': (r) => r.status === 200,
    'has payment reference': (r) => r.json('data.reference') !== undefined,
  });
  
  if (!success) {
    errorRate.add(1);
    return null;
  }
  
  return response.json('data');
}

function simulatePaymentCallback(reference) {
  const callbackData = {
    reference: reference,
    status: 'PAID',
    amount: 100000,
    paid_at: new Date().toISOString(),
    payment_method: 'BRIVA'
  };
  
  const response = http.post(`${BASE_URL}/v1/payments/callback`, JSON.stringify(callbackData), {
    headers: { 'Content-Type': 'application/json' },
  });
  
  check(response, {
    'payment callback processed': (r) => r.status === 200,
  }) || errorRate.add(1);
}