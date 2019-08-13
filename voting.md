## Voting with Redis

To upvote or downvote an article, we can use two keys, one to store the headline (`article:<id>:headline`), another to store the votes (`article:<id>:votes`).

```
SET article:12345:headline 'Headline 1'
SET article:10001:headline 'Headline 2'

INCR article:12345:votes
```

```js
const redis = require('redis')
const client = redis.createClient()

function upVote(id) {
  const key = `article:${id}:votes`
  client.incr(key)
}

function downVote(id) {
  const key = `article:${id}:votes`
  client.decr(key)
}

function showResults(id) {
  const headlineKey = `article:${id}:headline`
  const voteKey = `article:${id}:headline`
  client.mget([headlineKey, voteKey], (err, replies) => {
    console.log(`The article "${replies[0]}" has ${replies[1]} votes.`)
  })
}

upVote(12345) // Article 12345 has 1 vote
upVote(12345) // Article 12345 has 2 votes
upVote(12345) // Article 12345 has 3 votes
upVote(10001) // Article 10001 has 1 vote
downVote(12345) // Article 12345 has 2 votes

showResults(12345) // The article "headline 1" has 2 votes.
showResults(10001) // The article "headline 2" has 1 vote.
```
