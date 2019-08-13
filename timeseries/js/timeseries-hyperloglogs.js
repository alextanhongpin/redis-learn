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
      '1sec': {
        name: '1sec',
        ttl: 2 * time.Hour,
        duration: time.Second
      },
      '1min': {
        name: '1min',
        ttl: 7 * time.Day,
        duration: time.Minute
      },
      '1hour': {
        name: '1hour',
        ttl: 60 * time.Day,
        duration: time.Hour
      },
      '1day': {
        name: '1day',
        ttl: null,
        duration: time.Day
      }
    }
  }

  insert (timestampInSeconds, thing) {
    for (let precKey in this.precisions) {
      const prec = this.precisions[precKey]
      const key = this._getKeyName(prec, timestampInSeconds)
      this.client.pfadd(key, thing)
      if (prec.ttl !== null) {
        this.client.expire(key, prec.ttl)
      }
    }
  }

  _getKeyName (precision, timestampInSeconds) {
    const roundedTimestamp = this._getRoundedTimestamp(
      timestampInSeconds,
      precision.duration
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
      multi.pfcount(key)
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
