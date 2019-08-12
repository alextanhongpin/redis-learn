// http://www.tom-e-white.com/2007/11/consistent-hashing.html

const crypto = require('crypto')

const hashFunction = {
  hash (key) {
    const hash = crypto
      .createHash('md5')
      .update(key)
      .digest('hex')
    console.log(hash)
    return parseInt(hash, 16)
  }
}

console.log(hashFunction.hash('hello world'))

class ConsistentHashing {
  constructor (numberOfReplicas, hashFunction = hashFunction, nodes = []) {
    this.numberOfReplicas = numberOfReplicas
    this.hashFunction = hashFunction
    this.circle = {}

    for (let node of nodes) {
      this.add(node)
    }
  }
  add (node) {
    for (let i = 0; i < this.numberOfReplicas; i++) {
      const key = this.hashFunction.hash(node.toString() + ':' + i)
      this.circle[key] = node
    }
  }
  remove (node) {
    for (let i = 0; i < this.numberOfReplicas; i += 1) {
      const key = this.hashFunction.hash(node.toString() + ':' + i)
      delete this.circle[key]
    }
  }

  get (key) {
    const hash = this.hashFunction.hash(key)
    const keys = Object.keys(this.circle).sort()
    if (!keys.length) return null

    for (let key of keys) {
      if (key >= hash) {
        return this.circle[key]
      }
    }

    return this.circle[keys[0]]
  }
}

const consistentHashing = new ConsistentHashing(3, hashFunction, [
  'host-a',
  'host-b',
  'host-c'
])

console.log('initial', consistentHashing)
consistentHashing.remove('host-a')
console.log('host-a removed', consistentHashing)
console.log(consistentHashing.get('key1'))
console.log(consistentHashing.get('key2'))
console.log(consistentHashing.get('key3'))
