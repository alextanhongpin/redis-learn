const redis = require('redis')
const client = redis.createClient({ password: 123456 })

function transfer (from, to, value, callback) {
  client.get(from, (err, balance) => {
    if (err) return console.log(err)
    const multi = client.multi()
    multi.decrby(from, value)
    multi.incrby(to, value)
    if (balance >= value) {
      multi.exec((err, reply) => {
        if (err) return console.log(err)
        callback(null, reply[0])
      })
    } else {
      multi.discard()
      callback(new Error('insufficient funds'), null)
    }
  })
}
async function main () {
  client.flushall()
  client.mset('max:checkings', 100, 'hugo:checkings', 100, (err, reply) => {
    console.log('Max checkings: 100')
    console.log('Hugo checkings: 100')
    transfer('max:checkings', 'hugo:checkings', 40, (err, balance) => {
      if (err) return console.log(err)
      console.log('Transferred 40 from Max to Hugo')
      console.log('Max balance:', balance)
      client.quit()
    })
  })
}

main().catch(console.error)
