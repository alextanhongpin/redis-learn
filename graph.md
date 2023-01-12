# RedisGraph

## Basic example

Create a small graph representing a subset of motorcycle riders and teams taking part in the MotoGP championship.

```redis
graph.query MotoGP "create (:Rider {name: 'Valentino Rossi'})-[:rides]->(:Team {name: 'Yamaha'}), (:Rider {name: 'Dani Pedrosa'})-[:rides]->(:Team {name: 'Honda'}), (:Rider {name: 'Andrea Dovizioso'})-[:rides]->(:Team {name: 'Ducati'})"
```

After creating the MotoGP graph, we can start asking question.

1. Who's riding for Team Yamaha?
```redis
graph.query MotoGP "MATCH (r:Rider)-[:rides]->(t:Team) WHERE t.name = 'Yamaha' RETURN r.name, t.name"
```

1. How many riders represent team Ducati?
```redis
graph.query MotoGP "match (r:Rider)-[:rides]->(:Team {name: 'Ducati'}) RETURN count(r)"
```

```
-- This won't work if the user doesn't exists
graph.query MotoGP "match (r:Rider {name: 'John'}) create (r)-[:rides]->(:Team {name: 'Yamaha'}) RETURN r"
graph.query MotoGP "match (r:Rider {name: 'John'}) set r.age = 10"
graph.query MotoGP "match (r:Rider)-[:rides]->(:Team {name: 'Yamaha'}) RETURN r"

-- Use merge instead
graph.query MotoGP "MERGE (r:Rider {name: 'John'}) create (r)-[:rides]->(:Team {name: 'Yamaha'}) RETURN r"

-- We can add new property
graph.query MotoGP "match (t: Team{name: 'Yamaha'}) set t.year = 2022"
graph.query MotoGP "match (t: Team{name: 'Yamaha'}) RETURN t.year"
graph.query MotoGP "match (t: Team{name: 'Yamaha'}) set t.year = null"
```


## Social Media

```
alice
bob
charles

alice is friend with bob
bob is friend with charles
```

```
-- Create graph
graph.query social "create (:User {name: 'alice'})-[:FRIENDS_WITH]->(:User {name: 'bob'}), (:User {name: 'bob'})-[:FRIENDS_WITH]->(:User {name: 'charles'})"

-- Return friends of alice
graph.query social "match (:User {name: 'alice'})-[]->(fof) return fof"

graph.query social "match (:User {name: 'alice'})-[]->()<-[]-(fof) return fof"

graph.delete social
```

Scenario 2:
```
alice is friend with charles
bob is friend with charles
```

```
-- Create graph
graph.query social "create (:User {name: 'alice'})-[:FRIENDS_WITH]->(:User {name: 'charles'}), (:User {name: 'bob'})-[:FRIENDS_WITH]->(:User {name: 'charles'})"

-- Return friends of alice
graph.query social "match (:User {name: 'alice'})-[]->(fof) return fof"
graph.query social "match (:User {name: 'bob'})-[]->(fof) return fof"

-- This returns both alice and bob
graph.query social "match (:User {name: 'charles'})-[:FRIENDS_WITH]-(fw) return fw"

-- This returns none, since no direct relationship between charles with alice and bob
graph.query social "match (:User {name: 'charles'})-[:FRIENDS_WITH]->(fw) return fw"

-- Find users with the same friends.
graph.query social "match (:User)-[]->()<-[]-(fof) return fof"

-- Expected to get alice's friend's friend, which is charles's friends, bob.
-- But doesn't work. The same query seems to work with memgraph.
graph.query social "match (:User{name: 'alice'})-[]-()-[]-(fof) return fof"
graph.query social "match (:User{name: 'alice'})-[]-(fw)-[]-(fof) return fof"

-- This seems to work.
graph.query social "match (me:User{name: 'alice'})-[]->(fw) match (User{name: fw.name})<-[]-(fof) where me <> fof return fof"

-- Add two more nodes
damian is friend with alice
elie is friend with charles

-- It now becomes
alice is friend with charles
bob is friend with charles
damian is friend with alice
elie is friend with charles

graph.query social "create (:User {name: 'damian'})-[:FRIENDS_WITH]->(:User {name: 'alice'}), (:User {name: 'elie'})-[:FRIENDS_WITH]->(:User {name: 'charles'})"


-- Unidirectional friend with charles
graph.query social "match (me:User{name: 'alice'}) match (me)-[:FRIENDS_WITH]->(fw) return fw"

-- Bidirectional friends with damian and charles
graph.query social "match (me:User{name: 'alice'}) match (me)-[:FRIENDS_WITH]-(fw) return fw"

graph.query social "match (me:User{name: 'alice'}) match (me)-[:FRIENDS_WITH]-(fw)-[:FRIENDS_WITH]-(fof) return fof"

-- Charles is connected to alice, bob and elie
graph.query social "match (me:User{name: 'charles'})-[]-(fw) return fw.name"

-- But charles does not have direct relationship with them.
graph.query social "match (me:User{name: 'charles'})-[]->(fw) return fw.name"

-- alice knows charles, who knows bob and elie
graph.query social "match (me:User{name: 'alice'})-[]-(fw) match (:User{name: fw.name})-[]-(fof) where me <> fof return distinct(fof.name)"

-- This doesn't work somehow ...
graph.query social "match (me:User{name: 'alice'})-[:FRIENDS_WITH*2]-(fof) where me <> fof return distinct(fof.name)"
graph.query social "match (me:User{name: 'alice'})-[]->(fw)<-[]-(fof) where me <> fof return fw.name, fof.name"
graph.query social "match (me:User{name: 'alice'})--()--(fof:User) return fof.name"
graph.query social "match (me:User{name: 'alice'})-[]-(fw)-[:FRIENDS_WITH]->(fof) return fof"
graph.query social "match (me:User{name: 'alice'})-[:FRIENDS_WITH]-()-[:FRIENDS_WITH]-(fof) where me <> fof return fof"

graph.delete social
```
