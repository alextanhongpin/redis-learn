const redis = require('redis')
const client = redis.createClient({ password: 123456 })

function zpop (key, callback) {
  client.watch(key, (watchErr, watchReply) => {
    // if (watchErr) return console.log(watchErr)
    client.zrange(key, 0, 0, (zrangeErr, zrangeReply) => {
      if (zrangeErr) return console.log(zrangeErr)
      const multi = client.multi()
      multi.zrem(key, zrangeReply)
      multi.exec((transactionErr, transactionReply) => {
        if (transactionReply) {
          callback(zrangeReply[0])
        } else {
          zpop(key, callback)
        }
      })
    })
  })
}
async function main () {
  client.zadd('presidents', 1732, 'George Washington')
  client.zadd('presidents', 1809, 'Abraham Lincoln')
  client.zadd('presidents', 1858, 'Theodore Roosevelt')

  zpop('presidents', member => {
    console.log('The first president in the group is:', member)
    client.quit()
  })
}

main().catch(console.error)
