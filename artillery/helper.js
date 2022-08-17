function setRandomFlightId (context, events, done) {
  context.vars.flightId = context.vars.flights[Math.floor(Math.random() * context.vars.flights.length)].id
  return done()
}

function setBookingRequest (context, events, done) {
  const pickedSeats = []
  function pickRandomSeat () {
    for (let i = 0; i < 10; i++) {
      const seat = context.vars.seats[Math.floor(Math.random() * context.vars.seats.length)].seat
      if (!pickedSeats.includes(seat)) {
        pickedSeats.push(seat)
        return seat
      }
    }
    throw new Error('Could not find an available seat')
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

module.exports = {
  setRandomFlightId,
  setBookingRequest,
  setRandomDestination
}
