const url = require('url')

function setRandomFlightId (context, events, done) {
  context.vars.flightId = context.vars.flights[Math.floor(Math.random() * context.vars.flights.length)].id
  return done()
}

function setBookingRequest (context, events, done) {
  if (context.vars.seats.error === 'no seats available') {
    // no seats available
    // create dummy booking request that will fail
    context.vars.bookingRequest = {
      flightId: context.vars.flightId,
      passengers: [{}]
    }
    return done()
  }

  const pickedSeats = []
  function pickRandomSeat () {
    for (let i = 0; i < 3; i++) {
      const seat = context.vars.seats[Math.floor(Math.random() * context.vars.seats.length)].seat
      if (!pickedSeats.includes(seat)) {
        pickedSeats.push(seat)
        return seat
      }
    }
    console.error('Could not find an available seat')
    return null
  }

  const passengers = Array.from(Array(2))
    .map((_v, i) => ({
      name: `Passenger ${i}`,
      seat: pickRandomSeat()
    }))
  context.vars.bookingRequest = {
    flightId: context.vars.flightId,
    passengers
  }
  return done()
}

function setRandomDestination (context, events, done) {
  context.vars.flightFrom = context.vars.destinations.from[Math.floor(Math.random() * context.vars.destinations.from.length)]
  context.vars.flightTo = context.vars.destinations.to[Math.floor(Math.random() * context.vars.destinations.to.length)]
  return done()
}

// adapted from https://github.com/artilleryio/artillery-plugin-metrics-by-endpoint/blob/master/index.js
function afterResponse (requestParams, response, context, ee, next) {
  let basePath = (new url.URL(requestParams.url)).pathname
  if (basePath.endsWith('/seats')) {
    // remove flight id from request path
    const pathEl = basePath.split('/')
    pathEl[2] = '$flightID'
    basePath = pathEl.join('/')
  }
  const metricName = `scenario.${requestParams.scenarioName}.${requestParams.method}.${basePath}`
  ee.emit('counter', `${metricName}.response.${response.statusCode}`, 1)
  ee.emit('histogram', `${metricName}.total`, response.timings.phases.total)
  // ee.emit('histogram', `${metricName}.firstByte`, response.timings.phases.firstByte)
  // console.log(metricName, response.timings.phases)
  return next()
}

module.exports = {
  setRandomFlightId,
  setBookingRequest,
  setRandomDestination,
  afterResponse
}
