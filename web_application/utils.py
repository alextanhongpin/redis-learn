import redis

def connect():
    conn = redis.Redis(decode_responses=True)
    return conn