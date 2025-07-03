export const options = {
  stages: [
    // Tahap Warm-up: Memanaskan sistem dengan beban rendah
    { duration: '30s', target: 20 }, // Naik ke 20 VU dalam 30 detik

    // Tahap Ramp-up: Peningkatan bertahap menuju beban puncak yang diharapkan
    { duration: '1m', target: 50 },  // Naik ke 50 VU dalam 1 menit
    { duration: '1m', target: 100 }, // Naik ke 100 VU dalam 1 menit
    { duration: '1m', target: 200 }, // Naik ke 200 VU dalam 1 menit
    { duration: '1m', target: 300 }, // Naik ke 300 VU dalam 1 menit

    // Tahap Sustain: Pertahankan beban puncak selama periode yang cukup lama
    { duration: '5m', target: 250 }, // Pertahankan 250 VU selama 5 menit

    // Tahap Ramp-down: Turunkan beban secara bertahap
    { duration: '30s', target: 0 },  // Turun ke 0 VU dalam 30 detik
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