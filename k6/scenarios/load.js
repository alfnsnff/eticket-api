export const options = {
  stages: [
    { duration: '30s', target: 100 },   // pemanasan awal
    { duration: '30s', target: 200 },
    { duration: '30s', target: 300 },
    { duration: '30s', target: 400 },
    { duration: '30s', target: 500 },   // mencapai puncak
    { duration: '2m', target: 500 },    // pertahankan load puncak
    { duration: '1m', target: 0 },      // cool down
  ],
  
  thresholds: {
    // Toleransi kegagalan yang SANGAT KETAT untuk Load Test
    'http_req_failed': ['rate<0.05'], // Kurang dari 1% request yang gagal
    'http_req_duration': ['p(95)<3000'], // 95% request selesai dalam 1 detik

  },
  
  // Konfigurasi output tambahan
  summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(90)', 'p(95)', 'p(99)'],
  summaryTimeUnit: 'ms',
  
  // Konfigurasi untuk hasil yang lebih detail
  noConnectionReuse: false,
  userAgent: 'K6LoadTest/1.0',
};