export const options = {
  stages: [
    { duration: '1m', target: 50 },   // Naik ke 50 VU
    { duration: '1m', target: 100 },  // Naik ke 100 VU
    { duration: '1m', target: 200 },  // Naik ke 200 VU
    { duration: '1m', target: 300 },  // Naik ke 300 VU
    { duration: '1m', target: 400 },  // Naik ke 400 VU
    { duration: '3m', target: 500 },  // Bertahan di puncak VU
    { duration: '1m', target: 0 },    // Selesai
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'], // 95% request harus di bawah 1s (boleh gagal saat stress)
    http_req_failed: ['rate<0.1'],     // Maks 10% error rate
  },
};
