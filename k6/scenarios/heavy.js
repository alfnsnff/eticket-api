export const options = {
  stages: [
    { duration: '2m', target: 50 },    // pemanasan
    { duration: '3m', target: 250 },   // naik sedang
    { duration: '4m', target: 500 },   // naik berat
    { duration: '3m', target: 200 },   // turun bertahap
    { duration: '2m', target: 0 },     // ramp down
  ],
};
