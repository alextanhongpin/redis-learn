const crypto = require('crypto')

class HashPartitioning {
  constructor (clients) {
    this.clients = clients
  }
  get (key) {
    const index = this.hashFunction(key) % this.clients.length
    return this.clients[index]
  }
  hashFunction (str) {
    const hash = crypto.createHash('md5').update(str).digest('hex')
    // Convert hexadecimal to int.
    return parseInt(hash, 16)
  }
}

const hash = new HashPartitioning(['srv-1', 'srv-2'])
console.log(hash.get('john'))
console.log(hash.hashFunction('john'))

class ConsistentHashing {
  constructor (clients, vnodes) {
    this.vnodes = vnodes
    this.ring = {}
    this.setupRing(clients)
  }
  setupRing (clients) {
    for (const client of clients) {
      this.addClient(client)
    }
  }
  tagging (str) {
    // user:1:{users}
    // user:2:{users}
    const regex = /.+\{(.+)\}/
    const matches = regex.exec(str)
    if (!matches) {
      return str
    }
    return matches[1]
  }
  hashFunction (str) {
    const tag = this.tagging(str)
    return crypto.createHash('md5').update(tag).digest('hex')
  }
  addClient (client) {
    for (let i = 0; i < this.vnodes.length; i += 1) {
      const hash = this.hashFunction(client.addrss + ':' + i)
      this.ring[hash] = client
    }
  }
  get (key) {
    const ringHashes = Object.keys(this.ring)
    ringHashes.sort()
    const keyHash = this.hashFunction(key)
    for (const ringHash of ringHashes) {
      if (ringHash >= keyHash) {
        return this.ring[ringHash]
      }
    }
    // Fallback.
    return this.ring[ringHashes[0]]
  }
  delete (client) {
    for (let i = 0; i < this.vnodes.length; i += 1) {
      const hash = this.hashFunction(client.address + ':' + i)
      delete this.ring[hash]
    }
  }
}
