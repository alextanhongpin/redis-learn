class Queue {
  constructor (name, client) {
    this.name = name
    this.client = client
    this.queueKey = `queues:${name}`
    this.timeout = 0 // No timeout.
  }

  size (callback) {
    return this.client.llen(this.queueKey, callback)
  }

  // Items are inserted at the start of the list...
  push (data) {
    // Push the items to the start of the queue.
    this.client.lpush(this.queueKey, data)
  }

  // And removed at the end of the list. This is a FIFO queue.
  pop (callback) {
    this.client.brpop(this.queueKey, this.timeout, callback)
  }
}

module.exports = Queue
