const redis = require('redis')
const client = redis.createClient({ password: 123456 })

function saveLink (id, author, title, link) {
  client.hmset(
    `link:${id}`,
    'author',
    author,
    'title',
    title,
    'link',
    link,
    'score',
    0
  )
}

function upVote (id) {
  client.hincrby(`link:${id}`, 'score', 1)
}

function downVote (id) {
  client.hincrby(`link:${id}`, 'score', -1)
}

function showDetails (id) {
  client.hgetall(`link:${id}`, (err, replies) => {
    console.log('Title:', replies.title)
    console.log('Author:', replies.author)
    console.log('Link:', replies.link)
    console.log('Score:', replies.score)
  })
}

saveLink(
  123,
  'alextanhongpin',
  "Alex's Github Page",
  'https://github.com/alextanhongpin'
)
upVote(123)
upVote(123)
upVote(123)
upVote(123)
downVote(123)
showDetails(123)
client.quit()
