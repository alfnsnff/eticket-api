export const options = {
  stages: [
    // Tahap Warm-up: Memanaskan sistem dengan beban rendah
    { duration: '30s', target: 20 }, // Naik ke 20 VU dalam 30 detik

    // Tahap Ramp-up: Peningkatan bertahap menuju beban puncak yang diharapkan
    { duration: '1m', target: 50 },  // Naik ke 50 VU dalam 1 menit
    { duration: '1m', target: 100 }, // Naik ke 100 VU dalam 1 menit
    { duration: '1m', target: 200 }, // Naik ke 200 VU dalam 1 menit
    // Sesuaikan 'target' VU di sini dengan perkiraan BEBAN PUNCAK HARIAN Anda.
    // Misalnya, jika Anda memperkirakan 250 VU secara bersamaan.
    { duration: '1m', target: 300 }, // Naik ke 250 VU dalam 1 menit

    // Tahap Sustain: Pertahankan beban puncak selama periode yang cukup lama
    // Ini krusial untuk mengamati stabilitas jangka panjang dan potensi memory leak
    { duration: '5m', target: 250 }, // Pertahankan 250 VU selama 5 menit

    // Tahap Ramp-down: Turunkan beban secara bertahap
    { duration: '30s', target: 0 },  // Turun ke 0 VU dalam 30 detik
  ],
  thresholds: {
    // Toleransi kegagalan yang SANGAT KETAT untuk Load Test
    'http_req_failed': ['rate<0.01'], // Kurang dari 1% request yang gagal (idealnya 0%)
    // Waktu respons yang diharapkan untuk 95% request (sesuaikan dengan SLA/SLO Anda)
    'http_req_duration': ['p(95)<1000'], // 95% request selesai dalam 1 detik (1000ms)
    // Jika ada metrik bisnis spesifik, misal:
    // 'group_duration{group:::*/claim/lock}': ['p(95)<300'], // Waktu respons endpoint lock cepat
    // 'group_duration{group:::*/claim/entry}': ['p(95)<500'], // Waktu respons endpoint entry
    // 'group_duration{group:::*/payment/callback}': ['p(95)<500'], // Waktu respons endpoint payment
  },
};