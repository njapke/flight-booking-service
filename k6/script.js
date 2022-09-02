import http from 'k6/http'
import { check, sleep } from 'k6'
import { b64encode } from 'k6/encoding'

export const options = {
  noVUConnectionReuse: true,
  systemTags: ['iter', 'status', 'method', 'url', 'name', 'check', 'error', 'error_code', 'scenario', 'expected_response'],
  scenarios: {
    searchFlights: {
      executor: 'per-vu-iterations',
      exec: 'searchFlights',
      vus: 50,
      iterations: 100,
      maxDuration: '30m',

    },
    searchAndBookFlight: {
      executor: 'per-vu-iterations',
      exec: 'searchAndBookFlight',
      vus: 1,
      iterations: 10,
      maxDuration: '30m',
    }
  }
}

function selectRandomElement (data) {
  return data[Math.floor(Math.random() * data.length)]
}

function selectRandomUniqueElements (data, count) {
  const res = []
  for (let i = 0; i < count * 3; i++) {
    const el = selectRandomElement(data)
    if (res.includes(el)) continue
    res.push(el)
    if (res.length === count) return res
  }
  return []
}

export function searchFlights () {
  const endpoint = `http://${__ENV.TARGET}`
  const destinationRes = http.get(http.url`${endpoint}/destinations`)
  if (destinationRes.status !== 200) return
  const destination = JSON.parse(destinationRes.body)
  sleep(0.5)
  http.get(http.url`${endpoint}/flights?from=${selectRandomElement(destination.from)}`, { responseType: 'none' })
}

export function searchAndBookFlight () {
  const endpoint = `http://${__ENV.TARGET}`
  const destinationRes = http.get(http.url`${endpoint}/destinations`)
  if (destinationRes.status !== 200) return
  const destination = JSON.parse(destinationRes.body)
  sleep(1)

  const flightsRes = http.get(http.url`${endpoint}/flights?from=${selectRandomElement(destination.from)}`)
  if (flightsRes.status !== 200) return
  const flights = JSON.parse(flightsRes.body)
  const randomFlight = selectRandomElement(flights)
  sleep(1)

  const bookingRequest = {
    flightId: randomFlight.id,
    passengers: []
  }
  const seatsRes = http.get(http.url`${endpoint}/flights/${randomFlight.id}/seats`)
  if (seatsRes.status === 200) {
    const seats = JSON.parse(seatsRes.body)
    bookingRequest.passengers = selectRandomUniqueElements(seats, 2)
      .map((v, i) => ({ name: `Passenger ${i}`, seat: v.seat }))
  } else {
    // no seats available, create a booking request that will fail
    bookingRequest.passengers = [{ name: 'Passenger', seat: 'XX' }]
  }

  sleep(Math.floor(Math.random() * 3))
  const res = http.post(http.url`${endpoint}/bookings`, JSON.stringify(bookingRequest), {
    headers: { Authorization: `Basic ${b64encode('user:pw')}` },
    responseType: 'none'
  })
  check(res, {
    'successful booking': (r) => r.status === 200
  })
}
