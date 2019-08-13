const redis = require('redis')
const client = redis.createClient({ password: 123456 })

function addVisit (date, user) {
  const key = `visits:${date}`
  client.pfadd(key, user)
}

function count (dates) {
  const keys = []
  dates.forEach((date, index) => {
    keys.push(`visits:${date}`)
  })

  client.pfcount(keys, (err, reply) => {
    if (err) return console.log(err)
    console.log(`Dates ${dates.join(', ')} has ${reply} visits`)
  })
}

function aggregateDate (date) {
  const keys = [`visits:${date}`]
  for (let i = 0; i < 24; i += 1) {
    keys.push(`visits:${date}T${i}`)
  }
  client.pfmerge(keys, (err, reply) => {
    if (err) return console.log(err)
    console.log('Aggregated Date:', date)
  })
}

async function main () {
  const MAX_USERS = 200
  const TOTAL_VISITS = 1_000
  for (let i = 0; i < TOTAL_VISITS; i++) {
    const username = `user_${Math.floor(1 + Math.random() * MAX_USERS)}`
    const hour = Math.floor(Math.random() * 24)
    addVisit(`2019-01-01T${hour}`, username)
  }

  count(['2019-01-01T0'])
  count(['2019-01-01T5', '2019-01-01T6', '2019-01-01T7'])

  aggregateDate('2019-01-01')
  count(['2019-01-01'])
  client.quit()
}
main().catch(console.error)
