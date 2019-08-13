// Lua script is atomic, hence we can skip using WATCH.

const redis = require('redis')
const client = redis.createClient({ password: 123456 })

async function main () {
  client.zadd('presidents', 1732, 'George Washington')
  client.zadd('presidents', 1809, 'Abraham Lincoln')
  client.zadd('presidents', 1858, 'Theodore Roosevelt')

  const luaScript = `
    local elements = redis.call('ZRANGE', KEYS[1], 0, 0)
    redis.call('ZREM', KEYS[1], elements[1])
    return elements[1]
  `

  client.eval(luaScript, 1, 'presidents', (err, reply) => {
    if (err) return console.log(err)
    console.log('The first president in the group is:', reply)
    client.quit()
  })
}
main().catch(console.error)
