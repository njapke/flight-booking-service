# flight-booking-service

### GET /flights

```json
[
  {
    "id": "fea320f4-8f9a-4483-af65-bd49d6838a83",
    "from": "MUM",
    "to": "STV",
    "departure": "2022-07-05T23:23:51.37547748Z",
    "arrival": "2022-07-06T03:27:51.37547748Z",
    "status": "delayed"
  }
]
```

### GET /flights/{id}/seats

```json
[
  {
    "flightId": "7546127e-9924-43b9-aa53-961fd480d795",
    "seat": "6C",
    "row": 6,
    "price": 433,
    "available": true
  }
]
```

### POST /bookings

```json
{
  "id": "a39e5a34-0e15-4e1e-934d-b55a34610fb4",
  "userId": "user",
  "flightId": "fea320f4-8f9a-4483-af65-bd49d6838a83",
  "price": 37,
  "status": "confirmed",
  "passengers": [
    {
      "name": "Chris",
      "seat": "4C"
    }
  ]
}
```


### GET /bookings

```json
[
  {
    "id": "a39e5a34-0e15-4e1e-934d-b55a34610fb4",
    "userId": "user",
    "flightId": "fea320f4-8f9a-4483-af65-bd49d6838a83",
    "price": 37,
    "status": "confirmed",
    "passengers": [
      {
        "name": "Chris",
        "seat": "4C"
      }
    ]
  }
]
```

# Useful Commands

```bash

go test -run=^# -bench=. -cpuprofile=./bench.out ./...

curl http://localhost:3000/debug/pprof/profile?seconds=10 > app-bench.out

go tool pprof -http :9999 ./app-bench.out

go tool pprof -http :9999 -diff_base=./bench-slow.out ./bench-fast.out
```
