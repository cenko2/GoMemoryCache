# GoMemoryCache

To build and run in local docker, run the following under src

```
docker build .\ -t distributedcache
docker-compose up --scale distributedcache=2
```

The api can be reached from 
http://localhost:8080/cache/[Key]
