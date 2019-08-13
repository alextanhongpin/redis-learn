const redis = require('redis')
const client = redis.createClient({ password: 123456 })

class Leaderboard {
  constructor (key) {
    this.key = key
  }

  addUser (username, score) {
    client.zadd([this.key, score, username], (err, replies) => {
      if (err) {
        return console.log(err)
      }
      console.log(`User ${username} added to the leaderboard`)
    })
  }

  removeUser (username) {
    client.zrem(this.key, username, (err, replies) => {
      if (err) return console.log(err)
      console.log(`User ${username} removed successfully.`)
    })
  }

  getUserScoreAndRank (username) {
    const leaderboardKey = this.key
    client.zscore(leaderboardKey, username, (err, zscoreReply) => {
      if (err) return console.log(err)
      client.zrevrank(leaderboardKey, username, (err, zrevrankReply) => {
        if (err) return console.log(err)
        console.log(`Details of ${username}:`)
        console.log(`Score: ${zscoreReply} Rank: ${zrevrankReply + 1}`)
      })
    })
  }

  showTopUsers (n) {
    client.zrevrange([this.key, 0, n - 1, 'WITHSCORES'], (err, replies) => {
      if (err) return console.log(err)
      console.log(`Top ${n} users:`)
      for (let i = 0; i < replies.length; i += 2) {
        console.log(`User: ${replies[i]} Score: ${replies[i + 1]}`)
      }
    })
  }

  getUsersAroundUser (username, n, callback) {
    const leaderboardKey = this.key
    client.zrevrank(leaderboardKey, username, (err, zrevrankReply) => {
      if (err) return console.log(err)

      let startOffset = Math.floor(zrevrankReply - Math.floor(n / 2) + 1)
      if (startOffset < 0) startOffset = 0
      let endOffset = startOffset + n - 1

      client.zrevrange(
        [leaderboardKey, startOffset, endOffset, 'WITHSCORES'],
        (err, zrevrangeReply) => {
          if (err) return console.log(err)
          const users = []
          for (let i = 0, rank = 1; i < zrevrangeReply.length; i += 2, rank++) {
            const user = {
              rank: startOffset + rank,
              username: zrevrangeReply[i],
              score: zrevrangeReply[i + 1]
            }
            users.push(user)
          }
          callback(users)
        }
      )
    })
  }
}

async function main () {
  const leaderboard = new Leaderboard('game-score')
  leaderboard.addUser('A', 10)
  leaderboard.addUser('B', 20)
  leaderboard.addUser('C', 15)
  leaderboard.addUser('D', 5)
  leaderboard.addUser('E', 22)

  leaderboard.getUserScoreAndRank('A')
  leaderboard.showTopUsers(3)
  leaderboard.getUsersAroundUser('A', 5, users => {
    console.log(`Users around: A:`)
    users.forEach(user => {
      console.log(`#${user.rank} User: ${user.username} Score: ${user.score}`)
    })
    client.quit()
  })
}

main().catch(console.error)
