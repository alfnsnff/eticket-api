export const options = {
  stages: [
    { duration: '30s', target: 50 },
    { duration: '1m', target: 100 },
    { duration: '1m', target: 150 },
    { duration: '2m', target: 200 },
    { duration: '1m', target: 0 },
  ],
  
  thresholds: {
    // Toleransi kegagalan yang SANGAT KETAT untuk Load Test
    'http_req_failed': ['rate<0.05'], // Kurang dari 5% request yang gagal
    'http_req_duration': ['p(95)<3000'], // 95% request selesai dalam 1 detik

  },
  
  // Konfigurasi output tambahan
  summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(90)', 'p(95)', 'p(99)'],
  summaryTimeUnit: 'ms',
  
  // Konfigurasi untuk hasil yang lebih detail
  noConnectionReuse: false,
  userAgent: 'K6LoadTest/1.0',
};