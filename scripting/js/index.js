const redis = require('redis')
const client = redis.createClient({ password: 123456 })

async function main () {
  // EVAL script numkeys key [key...] arg [arg...]
  // script: the lua script as a stirng
  // numkeys: the number of Redis keys being passed as the parameters to the script.
  // key: The key name that will be available through the variable KEYS inside the script.
  // arg: An additional argument that will be available through the variable ARGV inside the script
  // NOTE: Lua index starts from 1, not 0 for arrays.

  client.set('mykey', 'myvalue')
  const luaScript = `
    return redis.call("GET", KEYS[1])
  `
  client.eval(luaScript, 1, 'mykey', (err, reply) => {
    if (err) return console.log(err)
    console.log(reply)
    client.quit()
  })
}

main().catch(console.error)
