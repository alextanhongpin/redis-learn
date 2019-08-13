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
      // Stores a maximum of 120 timestamp of 1 second each (2 minutes of data points).
      '1sec': {
        name: '1sec',
        ttl: 2 * time.Hour,
        duration: time.Second,
        quantity: 2 * time.Minute
      },

      // Stores a maximum of 120 timestamps of 1 minute each (2 hours
      // of data points).
      '1min': {
        name: '1min',
        ttl: 7 * time.Day,
        duration: time.Minute,
        quantity: 2 * time.Hour
      },
      // Stores a maximum of 120 timestamp of 1 hour each (5 days of
      // data points).
      '1hour': {
        name: '1hour',
        ttl: 60 * time.Day,
        duration: time.Hour,
        quantity: 5 * time.Day
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
    // configuration (zset-max-ziplist-entries default is 128), so
    // any number smaller than 128 makes this solution more
    // memory-efficient than using the String solution.
  }

  insert (timestampInSeconds, thing) {
    for (let precKey in this.precisions) {
      const prec = this.precisions[precKey]
      const key = this._getKeyName(prec, timestampInSeconds)
      const timestampScore = this._getRoundedTimestamp(
        timestampInSeconds,
        prec.duration
      )
      const member = `${timestampScore}:${thing}`
      this.client.zadd(key, timestampScore, member)
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
      multi.zcount(key, timestamp, timestamp)
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

  const concurrentPlays = new TimeSeries(client, 'concurrentplays')
  const beginTimestamp = 0
  concurrentPlays.insert(beginTimestamp, 'user:max')
  concurrentPlays.insert(beginTimestamp, 'user:max')
  concurrentPlays.insert(beginTimestamp + 1, 'user:hugo')
  concurrentPlays.insert(beginTimestamp + 1, 'user:renata')
  concurrentPlays.insert(beginTimestamp + 3, 'user:hugo')
  concurrentPlays.insert(beginTimestamp + 61, 'user:john')

  function displayResults (precisionName, results) {
    console.log(`Results from: ${precisionName}`)
    console.table(results)
  }

  concurrentPlays.fetch(
    '1sec',
    beginTimestamp,
    beginTimestamp + 4,
    displayResults
  )
  concurrentPlays.fetch(
    '1min',
    beginTimestamp,
    beginTimestamp + 120,
    displayResults
  )

  client.quit()
}

main().catch(console.error)
