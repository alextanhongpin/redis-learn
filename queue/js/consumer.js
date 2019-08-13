const redis = require('redis')
const client = redis.createClient({ password: 123456 })
const Queue = require('./Queue')
const logsQueue = new Queue('logs', client)

function logMessages () {
  logsQueue.pop((err, replies) => {
    if (err) {
      console.log(err)
      return
    }
    const queueName = replies[0]
    const message = replies[1]
    console.log(`[consumer]: ${message}`)
    logsQueue.size((err, size) => {
      if (err) {
        console.log(err)
        return
      }
      console.log(`${size} logs left`)
    })
    logMessages()
  })
}

logMessages()
