const redis = require('redis')
// Create a client that uses nodejs buffers rather than string. It is
// simpler to manipulate bytes with buffers than with JavaScript
// strings.
const client = redis.createClient({
  password: 123456,
  return_buffers: true
})

function storeDailyVisit (date, userId) {
  const key = `visits:daily:${date}`
  client.setbit(key, userId, 1, (err, reply) => {
    if (err) return console.log(err)
    console.log(`User: {userId} visited on ${date}`)
  })
}

function countVisits (date) {
  const key = `visits:daily:${date}`
  client.bitcount(key, (err, reply) => {
    if (err) return console.log(err)
    console.log(`${date} had ${reply} visits`)
  })
}

function showUserIdsFromVisits (date) {
  const key = `visits:daily:${date}`
  client.get(key, (err, bitmapValue) => {
    if (err) return console.log(err)
    const userIds = []
    const data = bitmapValue.toJSON().data
    data.forEach((byte, byteIndex) => {
      for (let bitIndex = 7; bitIndex >= 0; bitIndex--) {
        const visited = (byte >> bitIndex) & 1
        if (visited === 1) {
          const userId = byteIndex * 8 + (7 - bitIndex)
          userIds.push(userId)
        }
      }
    })
    console.log(`Users ${userIds.join(', ')} visited on ${date}`)
  })
}

async function main () {
  storeDailyVisit('2019-01-01', '1')
  storeDailyVisit('2019-01-01', '2')
  storeDailyVisit('2019-01-01', '10')
  storeDailyVisit('2019-01-02', '1')

  countVisits('2019-01-01')
  countVisits('2019-01-02')
  showUserIdsFromVisits('2019-01-01')
  client.quit()
}

main().catch(console.error)
