# go-geo-redis

A sample app that shows how to use Redis’s geospatial features with Go.

This repo contains a simple command-line app that shows how to [Redis’s commands for geospatial
data](https://redis.io/commands#geo) with Go and the [go-redis](https://github.com/go-redis/redis)
library.


## Run

The easiest way to run Redis is with Docker:

```
docker run -p 6379:6379 redis
```

Then run the sample app like this:

```
go run main.go add
go run main.go lookup Zurich
go run main.go find Geneva
```


## Code

Each of the `connect`, `add`, `lookup` and `find` functions is self-contained code that shows one
operation on the database.
