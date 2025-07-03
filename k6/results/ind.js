{
    "root_group": {
        "name": "",
        "path": "",
        "id": "d41d8cd98f00b204e9800998ecf8427e",
        "groups": {
            "Get Schedules": {
                "path": "::Get Schedules",
                "id": "1ff9489e4aecbab4130dffcc3e5cd03a",
                "groups": {},
                "checks": {
                        "load schedules success": {
                            "name": "load schedules success",
                            "path": "::Get Schedules::load schedules success",
                            "id": "d7930b0532faef923b00592c74568498",
                            "passes": 3483,
                            "fails": 1
                        }
                    },
                "name": "Get Schedules"
            },
            "Get Schedule by ID": {
                "path": "::Get Schedule by ID",
                "id": "a3cc98a0b6f4e50a23ff0d152a3caf5f",
                "groups": {},
                "checks": {
                        "load class success": {
                            "passes": 3439,
                            "fails": 0,
                            "name": "load class success",
                            "path": "::Get Schedule by ID::load class success",
                            "id": "1749ff3bd43cceb661ed4207b608ebd1"
                        }
                    },
                "name": "Get Schedule by ID"
            },
            "Claim Lock": {
                "name": "Claim Lock",
                "path": "::Claim Lock",
                "id": "17097546a3fc691b5c0f8d9ca6060e51",
                "groups": {},
                "checks": {
                        "claim lock success": {
                            "passes": 3384,
                            "fails": 2,
                            "name": "claim lock success",
                            "path": "::Claim Lock::claim lock success",
                            "id": "285d7793f111caa953d1c803e4a70a97"
                        }
                    }
            },
            "Claim Entry": {
                "name": "Claim Entry",
                "path": "::Claim Entry",
                "id": "7127c92ddd7507150af1be4973a1cc3e",
                "groups": {},
                "checks": {
                        "claim entry success": {
                            "id": "11cbf43e5ec4ddf6af6487213c5868f2",
                            "passes": 0,
                            "fails": 3332,
                            "name": "claim entry success",
                            "path": "::Claim Entry::claim entry success"
                        }
                    }
            },
            "Get Booking by Order ID": {
                "name": "Get Booking by Order ID",
                "path": "::Get Booking by Order ID",
                "id": "faf20d84412f46bc411ba244c0948efe",
                "groups": {},
                "checks": {
                        "get booking success": {
                            "name": "get booking success",
                            "path": "::Get Booking by Order ID::get booking success",
                            "id": "7c912f6fb073f95581b95bef5877d45c",
                            "passes": 3277,
                            "fails": 0
                        }
                    }
            },
            "Payment Callback": {
                "id": "7f0747d5b7ab2dc2964e62eb9d76d72f",
                "groups": {},
                "checks": {
                        "payment callback": {
                            "name": "payment callback",
                            "path": "::Payment Callback::payment callback",
                            "id": "1402e9ec9d57e38c861bce20004ffa8a",
                            "passes": 3226,
                            "fails": 0
                        }
                    },
                "name": "Payment Callback",
                "path": "::Payment Callback"
            }
        },
        "checks": {}
    },
    "metrics": {
        "http_req_duration": {
            "p(95)": 6927.50588999999,
            "p(99)": 12966.030854999995,
            "avg": 1812.8683514793527,
            "min": 28.1663,
            "med": 836.8969,
            "max": 30050.0404,
            "p(90)": 4663.588480000001,
            "thresholds": {
                "p(95)<3000": true
            }
        },
        "http_req_tls_handshaking": {
            "avg": 0,
            "min": 0,
            "med": 0,
            "max": 0,
            "p(90)": 0,
            "p(95)": 0,
            "p(99)": 0
        },
        "iterations": {
            "count": 3226,
            "rate": 10.96007814637641
        },
        "vus_max": {
            "value": 300,
            "min": 300,
            "max": 300
        },
        "http_req_waiting": {
            "p(90)": 4662.655020000001,
            "p(95)": 6926.966414999991,
            "p(99)": 12965.777756999996,
            "avg": 1812.4644943804644,
            "min": 28.1663,
            "med": 836.8199,
            "max": 30049.5461
        },
        "http_req_duration{expected_response:true}": {
            "p(95)": 6949.685519999993,
            "p(99)": 12968.564687999928,
            "avg": 1829.004094562444,
            "min": 28.1663,
            "med": 861.5236,
            "max": 29200.1046,
            "p(90)": 4661.6605800000025
        },
        "group_duration": {
            "med": 837.1575,
            "max": 30050.0404,
            "p(90)": 4664.60368,
            "p(95)": 6928.17619999999,
            "p(99)": 12966.030854999995,
            "avg": 1813.6120164316944,
            "min": 28.1663
        },
        "vus": {
            "max": 300,
            "value": 300,
            "min": 1
        },
        "http_req_receiving": {
            "p(95)": 0.9983849999999999,
            "p(99)": 1.1264699999999996,
            "avg": 0.37651867553613894,
            "min": 0,
            "med": 0,
            "max": 1406.6189,
            "p(90)": 0.8947700000000001
        },
        "http_req_failed": {
            "passes": 3335,
            "fails": 16809,
            "thresholds": {
                "rate<0.05": true
            },
            "value": 0.16555798252581413
        },
        "http_req_blocked": {
            "avg": 0.5129654189833199,
            "min": 0,
            "med": 0,
            "max": 133.1968,
            "p(90)": 0,
            "p(95)": 0,
            "p(99)": 31.433841
        },
        "iteration_duration": {
            "avg": 10118.353049813995,
            "min": 242.5013,
            "med": 7125.1161,
            "max": 57787.2493,
            "p(90)": 23351.75355,
            "p(95)": 29724.902000000002,
            "p(99)": 41624.903875
        },
        "http_req_sending": {
            "p(90)": 0,
            "p(95)": 0,
            "p(99)": 0.7752,
            "avg": 0.027338423351866543,
            "min": 0,
            "med": 0,
            "max": 4.5561
        },
        "checks": {
            "passes": 16809,
            "fails": 3335,
            "value": 0.8344420174741859
        },
        "http_reqs": {
            "count": 20144,
            "rate": 68.43763613781972
        },
        "data_received": {
            "count": 16590430,
            "rate": 56364.6649975163
        },
        "data_sent": {
            "count": 4619593,
            "rate": 15694.699406216194
        },
        "http_req_connecting": {
            "p(99)": 31.257029999999997,
            "avg": 0.505804924543288,
            "min": 0,
            "med": 0,
            "max": 133.1968,
            "p(90)": 0,
            "p(95)": 0
        }
    }
}