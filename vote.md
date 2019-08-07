As a User,
I want to vote for an article that I like,
in order to show my support for them.


Permutations of variables
- article validity (more than 1 week, less than one week) (state)
- user voted (yes/no) (action)

How do we calculate such combinations?

We have two variables, so we have 2 boxes. 
- For the first box, we fill the number of possible choices that can be taken, which is 2. 
- For the second box, we fill the number of possible choices that can be taken too, which is 2.

So 2 x 2 gives 4 possible combinations. Take a look at all-pairs testing.
We can think of them as a decision tree branching out too. 

Pruning
We do not need to test all combinations. We can choose to prune scenarios that are not necessary. In the variables above, if the article is less than a week old, we do not need to process them at all. This is more of a business rule, but we can choose to prune the tests cases if we know that it would no longer be handled.

From this, we have four possible combinations
1. article validity (more than 1 week) + user voted (yes)
2. article validity (more than 1 week) + user voted (no)
3. article validity (less than 1 week) + user voted (yes)
4. article validity (less than 1 week) + user voted (no)

Given an article with state 1 
And the article has 0 vote
And the article has  0 score
When user has not yet voted
And user vote for the article
Then the vote should be 1 and the score should be 432

Given an article with state 2
And the article has 1 vote
And the article has 432 score
When user has already voted
And user vote again
Then the vote and score should be unchanged
And the vote should be 1 
And the score should be 423

Given an article with state 3 and 4
Then the user can no longer vote for the article
```python
ONE_WEEK_IN_SECONDS = 86400 * 7
VOTE_SCORE = 432 # 86400 / 200

def vote_article(conn, user, article):
  # Articles more than a week old cannot be voted anymore.
  cutoff = time.time() - ONE_WEEK_IN_SECONDS
  if conn.zscore('time:', article) < cutoff:
    return
  
  # Equivalent to .split() in JavaScript, except that it returns
  # a tuple and also the separator.
  article_id = article.partition(':')[-1]

  # If the user has not voted for the article, increment the article score and vote count.
  if conn.sadd('voted:', article_id, user):
    conn.zincrby('score:', article, VOTE_SCRE)
    conn.hincrby(article, 'votes', 1)


# As a user, I want to create a new article.
def post_article(conn, user, title, link):
  # Increment the article id.
  article_id = str(conn.incr('article:'))

  # Vote is valid only for a week.
  voted = f'voted:{article_id}'
  conn.sadd(voted, user)
  conn.expire(voted, ONE_WEEK_IN_SECONDS)

  now = time.time()
  article = f'article:{article_id}'
  conn.hmset(article, {
    'title': title,
    'link': link,
    'poster': user,
    'time': now,
    'votes': 1
  })
  conn.zadd('score:', article, now + VOTE_SCORE)
  conn.zadd('time:', article, now)

  return article_id

ARTICLES_PER_PAGE = 25

# As a user, I want to get the top scoring articles or most recent articles.
def get_article(conn, page, order='score:'):
  start = (page - 1) * ARTICLES_PER_PAGE
  end = start + ARTICLES_PER_PAGE - 1
  
  # Sorted sets are ordered in ascending order - to get the top articles,
  # we get the ids by the score in descending order.
  ids = conn.zrevrange(order, start, end)
  articles = []
  for id in ids:
    article_data = conn.hgetall(id)
    article_data['id'] = id
    articles.append(article_data)

  return articles

# As a user, 
# I want to group articles by topics, 
# In order to see only relevant articles.
def add_remove_group(conn, article_id, to_add=[], to_remove=[]):
  article = f'article:{article_id}'
  for group in to_add:
    conn.sadd(f'group:{group}', article)

  for group in to_remove:
    conn.srem(f'group:{group}', article)

# As a user,
# I want to get the articles by topics,
# In order to view relevant articles.
def get_group_articles(conn, group, page, order='score:'):
  key = f'{order}{group}'
  if not conn.exists(key):
    # We can perform an intersection between a sorted set and set.
    # In order to get the articles from the group, and their score,
    # we intersect the article ids in the group, and their scores
    # and store them in another temporary data structure.
    conn.zinterstore(key,
                     [f'group:{group}', order],
                     aggregate='max') # Take the max score. Sets have a score of 1.
  
    # We only keep the grouped articles for 60 seconds,
    # after 60 seconds, we will create a new one with the
    # latest scores.
    conn.expire(key, 60)
  return get_article(conn, page, key)
```
