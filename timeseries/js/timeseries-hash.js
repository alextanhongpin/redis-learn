const redis = require('redis')
const client = redis.createClient({ password: 123456 })

const time = {
  Second: 1,
  Minute: 60,
  Hour: 60 * 60,
  Day: 60 * 60 * 24
}

class TimeSeries {
  constructor (client, namespace) {
    this.client = client
    this.namespace = namespace

    this.precisions = {
      // Stores a maximum of 300 timestamp of 1 second each (5 minutes of data points).
      '1sec': {
        name: '1sec',
        ttl: 2 * time.Hour,
        duration: time.Second,
        quantity: 5 * time.Minute
      },

      // Stores a maximum of 480 timestamps of 1 minute each (8 hours
      // of data points).
      '1min': {
        name: '1min',
        ttl: 7 * time.Day,
        duration: time.Minute,
        quantity: 8 * time.Hour
      },
      // Stores a maximum of 240 timestamp of 1 hour each (10 days of
      // data points).
      '1hour': {
        name: '1hour',
        ttl: 60 * time.Day,
        duration: time.Hour,
        quantity: 10 * time.Day
      },
      // Stores a maximum of 30 timestamps of 1 day each (30 days of data points).
      '1day': {
        name: '1day',
        ttl: null,
        duration: time.Day,
        quantity: 30 * time.Day
      }
    }
    // The quantity numbers are chosen based on the default Redis
    // configuration (hash-max-ziplist-entries is 512), so any number
    // smaller than 512 makes this solution more memory-efficient
    // than using the String solution.
  }

  insert (timestampInSeconds) {
    for (let precKey in this.precisions) {
      const prec = this.precisions[precKey]
      const key = this._getKeyName(prec, timestampInSeconds)
      const fieldName = this._getRoundedTimestamp(
        timestampInSeconds,
        prec.duration
      )
      this.client.hincrby(key, fieldName, 1)
      if (prec.ttl !== null) {
        this.client.expire(key, prec.ttl)
      }
    }
  }

  _getKeyName (precision, timestampInSeconds) {
    const roundedTimestamp = this._getRoundedTimestamp(
      timestampInSeconds,
      precision.quantity
    )
    return [this.namespace, precision.name, roundedTimestamp].join(':')
  }

  _getRoundedTimestamp (timestampInSeconds, precision) {
    // Alternatively, timestampInSeconds - timestampInSeconds % precision
    return Math.floor(timestampInSeconds / precision) * precision
  }

  fetch (precisionName, beginTimestamp, endTimestamp, callback) {
    const precision = this.precisions[precisionName]
    const begin = this._getRoundedTimestamp(beginTimestamp, precision.duration)
    const end = this._getRoundedTimestamp(endTimestamp, precision.duration)

    const multi = this.client.multi()
    for (
      let timestamp = begin;
      timestamp <= end;
      timestamp += precision.duration
    ) {
      const key = this._getKeyName(precision, timestamp)
      const fieldName = this._getRoundedTimestamp(timestamp, precision.duration)
      multi.hget(key, fieldName)
    }

    multi.exec((err, replies) => {
      if (err) return console.log(err)
      const results = []
      for (let i = 0; i < replies.length; i += 1) {
        const timestamp = beginTimestamp + i * precision.duration
        const value = parseInt(replies[i], 10) || 0
        results.push({ timestamp, value })
      }
      callback(precisionName, results)
    })
  }
}

async function main () {
  client.flushall()

  const item1Purchases = new TimeSeries(client, 'purchases:item1')
  const beginTimestamp = 0
  item1Purchases.insert(beginTimestamp)
  item1Purchases.insert(beginTimestamp + 1)
  item1Purchases.insert(beginTimestamp + 1)
  item1Purchases.insert(beginTimestamp + 3)
  item1Purchases.insert(beginTimestamp + 61)

  function displayResults (precisionName, results) {
    console.log(`Results from: ${precisionName}`)
    console.table(results)
  }

  item1Purchases.fetch(
    '1sec',
    beginTimestamp,
    beginTimestamp + 4,
    displayResults
  )
  item1Purchases.fetch(
    '1min',
    beginTimestamp,
    beginTimestamp + 120,
    displayResults
  )

  client.quit()
}

main().catch(console.error)
