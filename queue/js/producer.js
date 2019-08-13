const redis = require('redis')
const client = redis.createClient({ password: 123456 })
const Queue = require('./queue')
const logsQueue = new Queue('logs', client)

const MAX = 5
for (let i = 0; i < MAX; i += 1) {
  logsQueue.push(`Hello world #${i}`)
}
console.log(`created ${MAX} logs`)
client.quit()
